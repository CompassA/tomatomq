/*
 * @Author: Tomato
 * @Date: 2026-05-16 22:51:29
 */
package config

import (
	"fmt"

	"github.com/spf13/viper"

	tomatocfg "github.com/compassa/tomatomq/internal/pkg/config"
)

type Config struct {
	Server ServerConfig         `mapstructure:"server"` // broker服务配置
	Log    tomatocfg.LogConfig  `mapstructure:"log"`    // 日志配置
	Etcd   tomatocfg.EtcdConfig `mapstructure:"etcd"`   // etcd配置
}

type ServerConfig struct {
	Port       int    `mapstructure:"port"`       // broker服务端口
	Group      string `mapstructure:"group"`      // broker分组, topic的所有队列只会被单个分组下的broker处理
	BrokerName string `mapstructure:"brokerName"` // broker的名称, broker在分组中的唯一标识
}

func LoadConfig(env tomatocfg.Env) *Config {
	// 读取broker配置文件
	v := viper.New()
	v.SetConfigName(tomatocfg.AppBrokerConfigYaml)
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
	return &cfg
}
