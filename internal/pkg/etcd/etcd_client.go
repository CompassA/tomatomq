/*
 * @Author: Tomato
 * @Date: 2026-05-17 23:52:03
 */
package etcd

import (
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewClient(endpoint []string, dialTimoutSeconds int) *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: time.Duration(dialTimoutSeconds) * time.Second,
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect to ectd: %w", err))
	}
	return cli
}
