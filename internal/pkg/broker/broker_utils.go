/*
 * @Author: Tomato
 * @Date: 2026-06-02 23:30:29
 * @LastEditTime: 2026-06-03 00:51:21
 */
package broker

import (
	"fmt"
	"strings"

	"github.com/compassa/tomatomq/internal/pkg/constant"
)

func BuildEctdKey(brokerGroup, brokerName string) string {
	return fmt.Sprintf("%s/%s/%s", constant.ETCD_BROKER_PREFIX, brokerGroup, brokerName)
}

func BuildEctdValue(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func ParseEctdKey(key string) (brokerGroup, brokerName string) {
	// 不校验
	part := strings.Split(key, "/")
	brokerGroup = part[2]
	brokerName = part[3]
	return
}
