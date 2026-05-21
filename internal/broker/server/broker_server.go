/*
 * @Author: Tomato
 * @Date: 2026-05-10 21:44:47
 */
package broker

import (
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/compassa/tomatomq/internal/broker/config"
)

type Server struct {
	Sessions   []Session    // 所有客户端连接
	ListenAddr string       // 监听端口
	Ln         net.Listener // 示例
}

// 客户端连接
type Session struct {
	conn       net.Conn      // 网络连接
	done       chan struct{} // 网络断开标识
	Send       chan []byte   // broker server -> mq client 报文发送缓存
	LastActive int64         // 最后一次收到数据的时间戳
}

func NewServer(listenAddr string) (*Server, error) {
	if ln, err := net.Listen("tcp", listenAddr); err != nil {
		return nil, err
	} else {
		return &Server{
			Sessions:   []Session{},
			ListenAddr: listenAddr,
			Ln:         ln,
		}, nil
	}
}

func (s *Server) Serve() {
	for {
		conn, err := s.Ln.Accept()
		if err != nil {
			config.AppLogger.Error("accept erorr", slog.Any("error", err))
			return
		}
		config.AppLogger.Info("NewConnection",
			slog.String("local", conn.LocalAddr().String()),
			slog.String("remote", conn.RemoteAddr().String()))

		session := NewSession(conn, 128)
		s.Sessions = append(s.Sessions, *session)
		go session.ReadLoop()
		go session.WriteLoop()
	}
}

func NewSession(conn net.Conn, sendChanSize int) *Session {
	return &Session{
		conn:       conn,
		done:       make(chan struct{}),
		Send:       make(chan []byte, sendChanSize),
		LastActive: time.Now().UnixMilli(),
	}
}

func (s *Session) ReadLoop() {
	// TCP连接读取异常时, 关闭Session
	defer close(s.done)

	// 死循环读取TCP字节流
	buf := make([]byte, 4096)
	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				// 客户端主动关闭连接, 处理残余数据
				if n > 0 {
					handleNewBytes(buf[:n])
				}
				config.AppLogger.Info("client session closed", slog.String("remote", s.conn.RemoteAddr().String()))
				return
			} else {
				config.AppLogger.Error("cilent read error",
					slog.String("remote", s.conn.RemoteAddr().String()),
					slog.Any("error", err))
				return
			}
		}

		handleNewBytes(buf[:n])
	}
}

func (s *Session) WriteLoop() {
	// TCP写入异常 or 接收到session关闭标记, 关闭连接
	defer s.conn.Close()

	// 死循环监听报文发送缓存与session关闭标记的就绪态
	for {
		select {
		case msg := <-s.Send:
			if _, err := s.conn.Write(msg); err != nil {
				config.AppLogger.Error("cilent write error",
					slog.String("remote", s.conn.RemoteAddr().String()),
					slog.Any("error", err))
				return
			}
		case <-s.done:
			return
		}
	}
}

func handleNewBytes(bytes []byte) {
}
