/*
 * @Author: Tomato
 * @Date: 2026-05-10 21:29:53
 */
package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/compassa/tomatomq/internal/broker/config"
	"github.com/compassa/tomatomq/internal/broker/meta"

	tomatocfg "github.com/compassa/tomatomq/internal/pkg/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// 加载配置
	env := tomatocfg.FetchEnv()
	cfg := config.LoadConfig(env)

	// 初始化Etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Etcd.Endpoints,
		DialTimeout: time.Duration(cfg.Etcd.DialTimoutSeconds) * time.Second,
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect to ectd: %w", err))
	}
	defer cli.Close()

	ctx, cancel := context.WithCancel((context.Background()))
	defer cancel()

	// 元数据上报ectd
	meta.StartReport(ctx, &cfg.Server, cli)

	config.AppLogger.Info("broker started",
		slog.String("group", cfg.Server.Group),
		slog.String("brokerName", cfg.Server.BrokerName))
}
