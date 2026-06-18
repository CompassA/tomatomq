/*
 * @Author: Tomato
 * @Date: 2026-06-02 23:01:13
 * @LastEditTime: 2026-06-19 01:53:25
 */
package meta

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	brokerutil "github.com/compassa/tomatomq/internal/pkg/broker"
	tomatoconstant "github.com/compassa/tomatomq/internal/pkg/constant"
	"github.com/compassa/tomatomq/pkg/tomatolog"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type BrokerCacheRepo struct {
	cli   *clientv3.Client // etcd client
	cache atomic.Value     // broker_group -> broker_ip list
}

type BrokerMeta struct {
	Addr  string // broker网络信息, 格式: ip:port
	Name  string // broker名称
	Group string // brokerGroup
}

func NewRepo(cli *clientv3.Client) *BrokerCacheRepo {
	return &BrokerCacheRepo{
		cli: cli,
	}
}

func (r *BrokerCacheRepo) GetBrokerByGroup(brokerGroup string) []*BrokerMeta {
	cache := r.loadCache()

	res, ok := cache[brokerGroup]
	if !ok {
		return []*BrokerMeta{}
	}
	return res
}

func (r *BrokerCacheRepo) StartWatch(ctx context.Context) error {
	// 全量获取broker信息
	resp, err := r.cli.Get(ctx, tomatoconstant.EtcdBrokerPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	cache := map[string][]*BrokerMeta{}

	for _, kv := range resp.Kvs {
		k := string(kv.Key)
		v := string(kv.Value)

		group, name := brokerutil.ParseEctdKey(k)

		cache[group] = append(cache[group], &BrokerMeta{
			Name:  name,
			Group: group,
			Addr:  v,
		})
	}
	r.setCache(cache)

	// 开启增量监听
	go r.watchLoop(ctx)

	return nil
}

func (r *BrokerCacheRepo) watchLoop(ctx context.Context) {
	// 日志模块
	logger := slog.Default().With(tomatolog.LoggerNameKey, tomatoconstant.EtcdLogger)
	c := context.Background()

	// 外层循环, 处理watch重连
	for {
		// 注册watch,
		ch := r.cli.Watch(c, tomatoconstant.EtcdBrokerPrefix, clientv3.WithPrefix(), clientv3.WithProgressNotify())
		logger.Info("broker watch registered")

	inner:
		// 内层循环, 消费watch
		for {
			select {
			// 程序退出
			case <-ctx.Done():
				logger.Info("stop broker watch loop")
				return
			// 长时间没收到响应, 启动重连
			case <-time.After(15 * time.Minute):
				logger.Warn("no progress notify, try to reconnect")
				break inner
			case watch, ok := <-ch:
				// 通道关闭, 尝试重连
				if !ok {
					logger.Warn("watch channel closed, try to reconnect")
					break inner
				}

				// etcd服务端变更, 尝试重连
				if watch.Canceled {
					logger.Warn("watch canceled, try to reconnect")
					break inner
				}

				// 服务端心跳推送, 忽略
				if watch.IsProgressNotify() {
					logger.Info("receive watch progress notify")
					continue
				}

				// 消费变更
				r.applyEvents(watch.Events, logger)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func (r *BrokerCacheRepo) applyEvents(events []*clientv3.Event, logger *slog.Logger) {
	// 拷贝缓存
	cache := r.loadCache()
	newCache := make(map[string][]*BrokerMeta, len(cache))
	for k, src := range cache {
		dst := make([]*BrokerMeta, len(src))
		copy(dst, src)
		newCache[k] = dst
	}

	// 更新缓存
	for _, event := range events {
		k := string(event.Kv.Key)
		v := string(event.Kv.Value)
		t := event.Type
		logger.Info("receive watch event",
			slog.String("key", k),
			slog.String("value", v),
			slog.String("type", t.String()))

		group, name := brokerutil.ParseEctdKey(k)
		brokers, has := newCache[group]

		switch t {
		case mvccpb.PUT:
			if !has {
				brokers = []*BrokerMeta{}
				newCache[group] = brokers
			}

			// 检查broker是否在数组中存在, 存在则更新, 不存在则append
			pos := -1
			for i, broker := range brokers {
				if broker.Name == name {
					pos = i
					break
				}
			}
			if pos != -1 {
				brokers[pos] = &BrokerMeta{
					Name:  name,
					Addr:  v,
					Group: group,
				}
			} else {
				newCache[group] = append(brokers, &BrokerMeta{
					Name:  name,
					Addr:  v,
					Group: group,
				})
			}
		case mvccpb.DELETE:
			if !has {
				logger.Warn("broker group not exist in the cache", slog.String("brokerGroup", group))
				continue
			}
			pos := -1
			for i, broker := range brokers {
				if broker.Name == name {
					pos = i
					break
				}
			}
			if pos != -1 {
				newCache[group] = append(brokers[0:pos], brokers[pos+1:]...)
			}

		}
	}

	// 替换原子变量
	r.setCache(newCache)
}

func (r *BrokerCacheRepo) loadCache() map[string][]*BrokerMeta {
	return r.cache.Load().(map[string][]*BrokerMeta)
}

func (r *BrokerCacheRepo) setCache(cache map[string][]*BrokerMeta) {
	r.cache.Store(cache)
}
