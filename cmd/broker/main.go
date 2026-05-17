/*
 * @Author: Tomato
 * @Date: 2026-05-10 21:29:53
 */
package main

import (
	"context"
	"log/slog"

	"github.com/compassa/tomatomq/internal/broker"
	"github.com/compassa/tomatomq/internal/pkg/etcd"
)

func main() {
	// 加载配置
	env, cfg := broker.LoadConfig()

	// 初始化日志
	broker.LoadLogger(env, &cfg.Log)
	defer broker.SyncLogger()

	// 初始化Etcd
	cli := etcd.NewClient(cfg.Etcd.Endpoints, cfg.Etcd.DialTimoutSeconds)
	defer cli.Close()

	ctx, cancel := context.WithCancel((context.Background()))
	defer cancel()

	// 元数据上报ectd
	broker.StartReport(ctx, &cfg.Server, cli)

	broker.AppLogger.Info("broker started",
		slog.String("group", cfg.Server.Group),
		slog.String("brokerName", cfg.Server.BrokerName))
}
