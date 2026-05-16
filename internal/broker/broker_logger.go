/*
 * @Author: Tomato
 * @Date: 2026-05-16 22:16:38
 */
package broker

import (
	"fmt"
	"log/slog"
	"os"

	logger "github.com/compassa/tomatomq/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var AppLogger *slog.Logger

// file config
func LoadLogger(env Env, logcfg *LogConfig) {
	if len(logcfg.Loggers) == 0 || len(logcfg.Writers) == 0 {
		panic("logger config empty")
	}

	writerMap := map[string]zapcore.Core{}
	for _, writercfg := range logcfg.Writers {
		level := fetchLevel(writercfg.Level)
		var core zapcore.Core
		switch writercfg.Type {
		case AppBrokerWriterFileType:
			core = logger.NewZapCore(*writercfg.Path, *writercfg.MaxSize, *writercfg.MaxBackups, *writercfg.MaxAge, level)
		case AppBrokerWriterStdoutType:
			core = logger.NewStdoutCore(level)
		case AppBrokerWriterStderrType:
			core = logger.NewStderrCore(level)
		}
		writerMap[writercfg.Id] = core
	}

	loggerMap := map[string]*slog.Logger{}
	for _, loggercfg := range logcfg.Loggers {
		cores := []zapcore.Core{}
		if len(loggercfg.Writers) == 0 {
			panic(fmt.Errorf("logger %s without any writer", loggercfg.Name))
		}
		for _, writerId := range loggercfg.Writers {
			core, ok := writerMap[writerId]
			if !ok {
				panic(fmt.Errorf("writter %s not found", writerId))
			}
			cores = append(cores, core)
		}
		loggerMap[loggercfg.Name] = slog.New(logger.NewZapHandler(cores))
	}

	AppLogger = fetchLogger("app-logger", env, loggerMap)
}

func fetchLogger(name string, env Env, loggerMap map[string]*slog.Logger) *slog.Logger {
	logger, ok := loggerMap[name]
	if ok {
		return logger
	}

	if env != AppBrokerEnvProd {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	panic(fmt.Errorf("logger %s not configured", name))
}

func fetchLevel(l WriterLevel) zapcore.Level {
	switch l {
	case AppBrokerLoggerDebug:
		return zap.DebugLevel
	case AppBrokerLoggerInfo:
		return zap.InfoLevel
	case AppBrokerLoggerWarn:
		return zap.WarnLevel
	case AppBrokerLoggerError:
		return zap.ErrorLevel
	case AppBrokerLoggerDPanic:
		return zap.DPanicLevel
	case AppBrokerLoggerPanic:
		return zap.PanicLevel
	case AppBrokerLoggerFatal:
		return zap.FatalLevel
	default:
		panic(fmt.Errorf("unknown log level %s", l))
	}
}
