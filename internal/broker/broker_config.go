/*
 * @Author: Tomato
 * @Date: 2026-05-16 22:51:29
 */
package broker

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	// 系统相关环境变量
	AppBrokerEnvironmentVarAppEnv = "APP_ENV" // 应用环境

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
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type LogConfig struct {
	Writers []WriterConfig `mapstructure:"writers"`
	Loggers []LoggerConfig `mapstructure:"loggers"`
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

func LoadConfig() (Env, *Config) {
	env := fetchEnv()

	v := viper.New()
	v.SetConfigName("broker-" + string(env) + ".yaml")
	v.SetConfigType("yaml")
	v.AddConfigPath("../../config") // 配置文件目录
	v.AddConfigPath(".")            // 兼容不同运行目录
	v.AutomaticEnv()                // 允许环境变量覆盖

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read broker config failed, %w", err))
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("unmarshal broker config failed, %w", err))
	}
	return env, &cfg
}

func fetchEnv() Env {
	env := os.Getenv(AppBrokerEnvironmentVarAppEnv)
	if env == "" {
		return AppBrokerEnvDev
	}
	return Env(env)
}
