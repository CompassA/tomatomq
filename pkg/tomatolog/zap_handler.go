/*
 * @Author: Tomato
 * @Date: 2026-05-17 22:05:37
 * @LastEditTime: 2026-06-17 01:14:43
 * 日志模块封装思路:
 * - 编写slog handler, 日志API使用slog, 底层通过ZapLogger输出日志
 * - 应用只初始化一个logger, 每次受理网络请求时, 基于rootlogger调用with方法绑定通用参数(如traceId), 派生请求logger, 后续的业务逻辑中通过FromCtx使用派生的logger
 * - handler中保存多个zaplogger, 运行时通过设置context指定将使用哪个zaplogger写日志
 */
package tomatolog

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 一些参数key
type (
	loggerNameKeyType struct{}
	ctxLoggerKeyType  struct{}
)

var (
	CtxLoggerKey  = ctxLoggerKeyType{}  // ctx中的logger
	LoggerNameKey = loggerNameKeyType{} // ctx中配置的模块日志参数, 对应ZapHandler.loggers的key
)

type ZapHandler struct {
	loggers       map[string]*zap.Logger // 配置文件中的logger名称 -> logger实现类
	defaultLogger *zap.Logger            // 默认的logger
	attrs         []zap.Field            // attrs
}

func NewZapLogger(cores []zapcore.Core) *zap.Logger {
	var core zapcore.Core
	corelen := len(cores)
	if corelen <= 0 {
		return nil
	} else if corelen == 1 {
		core = cores[0]
	} else {
		core = zapcore.NewTee(cores...)
	}
	return zap.New(core)
}

func NewZapHandler(loggers map[string]*zap.Logger, defaultLoggerName string) *ZapHandler {
	if defaultLogger, ok := loggers[defaultLoggerName]; ok {
		return &ZapHandler{
			loggers:       loggers,
			defaultLogger: defaultLogger,
		}
	}
	panic("default logger " + defaultLoggerName + " not found")
}

func FromCtx(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(CtxLoggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

func (h *ZapHandler) Handle(c context.Context, r slog.Record) error {
	// 添加属性参数和日志参数
	attrlen := len(h.attrs)
	fields := make([]zap.Field, 0, r.NumAttrs()+attrlen)
	if attrlen > 0 {
		for _, attr := range h.attrs {
			fields = append(fields, attr)
		}
	}
	r.Attrs(func(attr slog.Attr) bool {
		fields = append(fields, slogAttrToZapField(attr))
		return true
	})

	// 选择目标zap logger
	logger := h.defaultLogger
	if loggerName, ok := c.Value(LoggerNameKey).(string); ok {
		if l, ok := h.loggers[loggerName]; ok {
			logger = l
		}
	}

	// 输出日志
	logger.Log(convertLevel(r.Level), r.Message, fields...)
	return nil
}

func (h *ZapHandler) Enabled(_ context.Context, l slog.Level) bool {
	return true
}

func (h *ZapHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// 拷贝旧的属性
	oldAttrLen := len(h.attrs)
	newAttrs := make([]zap.Field, oldAttrLen, oldAttrLen+len(attrs))
	if oldAttrLen > 0 {
		copy(newAttrs, h.attrs)
	}

	// 添加新的属性
	for _, attr := range attrs {
		newAttrs = append(newAttrs, slogAttrToZapField(attr))
	}

	// 构造新的logger
	return &ZapHandler{
		loggers:       h.loggers,
		defaultLogger: h.defaultLogger,
		attrs:         newAttrs,
	}
}

func (h *ZapHandler) WithGroup(name string) slog.Handler {
	panic("not implemented")
}

func NewZapCore(path string, maxSize int, maxBackups int, maxAge int, level zapcore.Level) zapcore.Core {
	// encoder
	encoder := newEncoder()

	// writer
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,       // 文件位置
		MaxSize:    maxSize,    // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     maxAge,     // 保留旧文件的最大天数
		MaxBackups: maxBackups, // 保留旧文件的最大个数
		Compress:   false,      // 是否压缩/归档旧文件
	}
	writer := zapcore.AddSync(lumberJackLogger)

	// zap core
	return zapcore.NewCore(encoder, writer, level)
}

func NewStdoutCore(level zapcore.Level) zapcore.Core {
	encoder := newEncoder()
	return zapcore.NewCore(encoder, os.Stdout, level)
}

func NewStderrCore(level zapcore.Level) zapcore.Core {
	encoder := newEncoder()
	return zapcore.NewCore(encoder, os.Stderr, level)
}

func newEncoder() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 序列化时间。eg: 2022-09-01T19:11:35.921+0800
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 将Level序列化为全大写字符串。例如，将info level序列化为INFO。
	encoder := zapcore.NewJSONEncoder(encodeConfig)
	return encoder
}

func slogAttrToZapField(a slog.Attr) zapcore.Field {
	key := a.Key
	value := a.Value

	switch value.Kind() {
	case slog.KindBool:
		return zap.Bool(key, value.Bool())
	case slog.KindDuration:
		return zap.Duration(key, value.Duration())
	case slog.KindFloat64:
		return zap.Float64(key, value.Float64())
	case slog.KindInt64:
		return zap.Int64(key, value.Int64())
	case slog.KindString:
		return zap.String(key, value.String())
	case slog.KindTime:
		return zap.Time(key, value.Time())
	case slog.KindUint64:
		return zap.Uint64(key, value.Uint64())
	case slog.KindAny:
		return zap.Any(key, value.Any())
	case slog.KindGroup:
		panic("not implemented")
	default:
		panic("not implemented")
	}
}

func convertLevel(level slog.Level) zapcore.Level {
	switch {
	case level >= slog.LevelError:
		return zapcore.ErrorLevel
	case level >= slog.LevelWarn:
		return zapcore.WarnLevel
	case level >= slog.LevelInfo:
		return zapcore.InfoLevel
	default:
		return zapcore.DebugLevel
	}
}
