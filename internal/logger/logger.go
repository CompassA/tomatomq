package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var BrokerApp *zap.Logger

func InitBroker() error {
	if logger, err := newLogger("./broker/app-default.log", 512, 10, 3, zapcore.ErrorLevel); err != nil {
		return err
	} else {
		BrokerApp = logger
	}

	return nil
}

func Sync() {
	if BrokerApp != nil {
		BrokerApp.Sync()
	}
}

func newLogger(path string, maxSize int, maxBackups int, maxAge int, level zapcore.Level) (*zap.Logger, error) {
	encode := newEncodeConfig()
	writer, err := newWriterConfig(path, maxSize, maxBackups, maxAge)
	if err != nil {
		return nil, err
	}
	return zap.New(
		zapcore.NewCore(encode, writer, level),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	), nil
}

func newEncodeConfig() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 序列化时间。eg: 2022-09-01T19:11:35.921+0800
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 将Level序列化为全大写字符串。例如，将info level序列化为INFO。
	return zapcore.NewJSONEncoder(encodeConfig)
}

func newWriterConfig(path string, maxSize int, maxBackups int, maxAge int) (zapcore.WriteSyncer, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,       // 文件位置
		MaxSize:    maxSize,    // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     maxAge,     // 保留旧文件的最大天数
		MaxBackups: maxBackups, // 保留旧文件的最大个数
		Compress:   false,      // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger), nil
}
