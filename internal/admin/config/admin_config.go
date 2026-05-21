/*
 * @Author: Tomato
 * @Date: 2026-05-20 23:41:14
 */

package config

import (
	"fmt"
	"log/slog"

	tomatocfg "github.com/compassa/tomatomq/internal/pkg/config"
	"github.com/spf13/viper"
)

var AppLogger *slog.Logger

type Config struct {
	Database tomatocfg.DBConfig   `mapstructure:"database"` // broker服务配置
	Etcd     tomatocfg.EtcdConfig `mapstructure:"etcd"`     // etcd配置
	Log      tomatocfg.LogConfig  `mapstructure:"log"`      // 日志配置
}

func LoadConfig(env tomatocfg.Env) *Config {
	v := viper.New()
	v.SetConfigName("admin-" + string(env) + ".yaml")
	v.SetConfigType("yaml")
	v.AddConfigPath("../../config") // 配置文件目录
	v.AddConfigPath(".")            // 兼容不同运行目录
	v.AutomaticEnv()                // 允许环境变量覆盖

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read admin config failed, %w", err))
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("unmarshal admin config failed, %w", err))
	}

	loggerMap := tomatocfg.LoadLogger(env, cfg.Log)
	AppLogger = tomatocfg.FetchLogger("app-logger", env, loggerMap)

	return &cfg
}

func SyncLogger() {
	tomatocfg.SyncLogger([]*slog.Logger{AppLogger})
}
