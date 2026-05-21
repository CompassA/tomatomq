/*
 * @Author: Tomato
 * @Date: 2026-05-18 00:02:41
 */
package meta

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	config "github.com/compassa/tomatomq/internal/broker/config"
	tomatocfg "github.com/compassa/tomatomq/internal/pkg/config"
	tomatoconstant "github.com/compassa/tomatomq/internal/pkg/constant"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func StartReport(ctx context.Context, cfg *config.ServerConfig, cli *clientv3.Client) {
	// 创建租约
	leaseResp, err := cli.Grant(ctx, 30)
	if err != nil {
		panic(fmt.Errorf("create release failed: %w", err))
	}

	config.AppLogger.Info("create lease succeed",
		slog.String("mark", "EtcdKeepAlive"),
		slog.Int64("leaseId", int64(leaseResp.ID)),
		slog.Int64("ttl", leaseResp.TTL))

	// 写入broker信息
	ip := os.Getenv(tomatocfg.AppBrokerEnvironmentVarPodIp)
	key := fmt.Sprintf("%s/%s/%s", tomatoconstant.ETCD_BROKER_PREFIX, cfg.Group, cfg.BrokerName)
	value := fmt.Sprintf("%s:%d", ip, cfg.Port)
	_, err = cli.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		panic(fmt.Errorf("regiter broker meta failed: %w", err))
	}
	config.AppLogger.Info("report broker meta succeed",
		slog.String("mark", "EtcdKeepAlive"),
		slog.String("key", key),
		slog.String("value", value))

	// 定期续期
	keepAliveCh, err := cli.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		panic(fmt.Errorf("broker meta keep alive failed: %w", err))
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				config.AppLogger.Info("context is done, stop lease",
					slog.String("mark", "EtcdKeepAlive"),
					slog.Int64("leaseId", int64(leaseResp.ID)))
				return
			case resp := <-keepAliveCh:
				if resp == nil {
					config.AppLogger.Info("lease response is empty, restart report", slog.String("mark", "EtcdKeepAlive"))
					StartReport(ctx, cfg, cli)
					return
				}
				config.AppLogger.Info("lease keep-alice succeed",
					slog.String("mark", "EtcdKeepAlive"),
					slog.Int64("leaseId", int64(resp.ID)),
					slog.Int64("ttl", resp.TTL),
				)
			}
		}
	}()
}
