/*
 * @Author: Tomato
 * @Date: 2026-05-18 00:02:41
 */
package broker

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	tomatoconstant "github.com/compassa/tomatomq/internal/pkg/constant"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// type BrokerMetaManager struct {
// 	LeaseId clientv3.LeaseID
// 	Cli     *clientv3.Client
// 	Cfg     *ServerConfig
// }

func StartReport(ctx context.Context, cfg *ServerConfig, cli *clientv3.Client) {
	// 创建租约
	leaseResp, err := cli.Grant(ctx, 30)
	if err != nil {
		panic(fmt.Errorf("create release failed: %w", err))
	}
	AppLogger.Info("create lease succeed",
		slog.String("mark", "EtcdKeepAlive"),
		slog.Int64("leaseId", int64(leaseResp.ID)),
		slog.Int64("ttl", leaseResp.TTL))

	// 写入broker信息
	ip := os.Getenv(AppBrokerEnvironmentVarPodIp)
	key := fmt.Sprintf("%s/%s/%s", tomatoconstant.ETCD_BROKER_PREFIX, cfg.Group, cfg.BrokerName)
	value := fmt.Sprintf("%s:%d", ip, cfg.Port)
	_, err = cli.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		panic(fmt.Errorf("regiter broker meta failed: %w", err))
	}
	AppLogger.Info("report broker meta succeed",
		slog.String("mark", "EtcdKeepAlive"),
		slog.String("key", key),
		slog.String("value", value))

	// 定期续期
	keepAliveCh, err := cli.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		panic(fmt.Errorf("broker meta keep alive failed: %w", err))
	}

	go func() {
		for resp := range keepAliveCh {
			if resp == nil {
				AppLogger.Info("lease response is empty, restart report", slog.String("mark", "EtcdKeepAlive"))
				StartReport(ctx, cfg, cli)
				return
			}
			AppLogger.Info("lease keep-alice succeed",
				slog.String("mark", "EtcdKeepAlive"),
				slog.Int64("leaseId", int64(resp.ID)),
				slog.Int64("ttl", resp.TTL),
			)
		}
	}()
}
