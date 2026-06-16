/*
 * @Author: Tomato
 * @Date: 2026-05-20 23:33:29
 */
package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/compassa/tomatomq/pkg/tomatolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// 系统相关环境变量
	AppBrokerEnvironmentVarAppEnv = "APP_ENV" // 应用环境
	AppBrokerEnvironmentVarPodIp  = "POD_IP"  // IP环境变量, 设置了这个环境变量后, broker会直接将这个IP注册当作自己的机器ip, 注册至etcd

	// 环境标识
	AppBrokerEnvDev  Env = "dev"
	AppBrokerEnvProd Env = "prod"

	// 日志文件类型
	AppBrokerWriterFileType   WriterType = "FILE"
	AppBrokerWriterStdoutType WriterType = "STDOUT"
	AppBrokerWriterStderrType WriterType = "STDERR"

	// 日志等级
	AppBrokerLoggerDebug  WriterLevel = "DEBUG"
	AppBrokerLoggerInfo   WriterLevel = "INFO"
	AppBrokerLoggerWarn   WriterLevel = "WARN"
	AppBrokerLoggerError  WriterLevel = "ERROR"
	AppBrokerLoggerDPanic WriterLevel = "DPANIC"
	AppBrokerLoggerPanic  WriterLevel = "PANIC"
	AppBrokerLoggerFatal  WriterLevel = "FATAL"

	// broker配置文件, 读取二进制可执行文件同文件夹下的"broker-config.yaml"
	AppBrokerConfigYaml = "broker-config.yaml"
)

type DBConfig struct {
	Dsn                string `mapstructure:"dsn"`            // GORM连接
	MaxIdleConns       int    `mapstructure:"maxIdle"`        // 连接池最大空闲
	MaxOpenConns       int    `mapstructure:"maxOpen"`        // 连接池最大连接数
	ConnMaxLifeMinutes int    `mapstructure:"connMaxMinutes"` // 连接最长存活时间
}

type EtcdConfig struct {
	Endpoints         []string `mapstructure:"endpoints"`
	DialTimoutSeconds int      `mapstructure:"dialTimoutSeconds"`
}

type LogConfig struct {
	Writers           []WriterConfig `mapstructure:"writers"`           // 日志文件配置
	Loggers           []LoggerConfig `mapstructure:"loggers"`           // 日志实例配置, 一个日志实例组合多个日志文件
	DefaultLoggerName string         `mapstructure:"defaultLoggerName"` // 默认的日志实例
}

type WriterConfig struct {
	Id         string      `mapstructure:"id"`         // 文件配置id
	Type       WriterType  `mapstructure:"type"`       // 文件类型
	Level      WriterLevel `mapstructure:"level"`      // 日志等级
	Path       *string     `mapstructure:"path"`       // 文件路径
	MaxSize    *int        `mapstructure:"maxSize"`    // 文件最大大小 MB
	MaxBackups *int        `mapstructure:"maxBackups"` // 最多保留几个文件
	MaxAge     *int        `mapstructure:"maxAge"`     // 单日志文件最多保留几天
}

type LoggerConfig struct {
	Name    string   `mapstructure:"name"`    // logger id
	Writers []string `mapstructure:"writers"` // 关联的写入配置
}

type (
	WriterType  string
	WriterLevel string
	Env         string
)

func LoadLogger(env Env, logcfg *LogConfig) *slog.Logger {
	if len(logcfg.Loggers) == 0 || len(logcfg.Writers) == 0 {
		panic("logger config empty")
	}

	// 日志文件配置
	writerMap := map[string]zapcore.Core{}
	for _, writercfg := range logcfg.Writers {
		level := fetchLevel(writercfg.Level)
		var core zapcore.Core
		switch writercfg.Type {
		case AppBrokerWriterFileType:
			core = tomatolog.NewZapCore(*writercfg.Path, *writercfg.MaxSize, *writercfg.MaxBackups, *writercfg.MaxAge, level)
		case AppBrokerWriterStdoutType:
			core = tomatolog.NewStdoutCore(level)
		case AppBrokerWriterStderrType:
			core = tomatolog.NewStderrCore(level)
		}
		writerMap[writercfg.Id] = core
	}

	// 不同业务模块组合日志文件
	loggerMap := map[string]*zap.Logger{}
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
		loggerMap[loggercfg.Name] = tomatolog.NewZapLogger(cores)
	}

	// 构造rootLogger
	rootLogger := slog.New(tomatolog.NewZapHandler(loggerMap, logcfg.DefaultLoggerName))
	slog.SetDefault(rootLogger)
	return rootLogger
}

func FetchEnv() Env {
	env := os.Getenv(AppBrokerEnvironmentVarAppEnv)
	if env == "" {
		return AppBrokerEnvDev
	}
	return Env(env)
}

func FetchPodIp() string {
	return os.Getenv(AppBrokerEnvironmentVarPodIp)
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
