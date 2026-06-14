/*
 * @Author: Tomato
 * @Date: 2026-05-27 23:12:30
 * @LastEditTime: 2026-06-14 22:09:20
 */
package mqadmin

import (
	"fmt"

	"github.com/compassa/tomatomq/internal/pkg/model"
)

type DatabaseRegisterReq struct {
	BrokerGroup string `json:"brokerGroup" binding:"required"`
	Name        string `json:"name" binding:"required"`
	User        string `json:"user" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Host        string `json:"host" binding:"required"`
	Port        int    `json:"port" binding:"gte=1025,lte=65535"`
}

func buildNewDatabase(req *DatabaseRegisterReq) *model.BrokerGroupDatabase {
	return &model.BrokerGroupDatabase{
		BrokerGroup: req.BrokerGroup,
		Dsn:         fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", req.User, req.Password, req.Host, req.Port, req.Name),
		Guid:        fmt.Sprintf("%s:%s", req.BrokerGroup, req.Name),
	}
}

type TopicCreateMode int

const (
	UnknownTopicCreateMode TopicCreateMode = iota
	// 普通模式, 每个Broker可管理一个topic的一个分区
	Simple
	// 严格模式, 每个Broker仅可管理Topic的一个分区, 当broker数量少于Topic分区数量时, Topic无法创建
	Restrict
)

type TopicRegisterReq struct {
	TopicName   string          `json:"topicName" binding:"required"`
	BrokerGroup string          `json:"brokerGroup" binding:"required"`
	Type        model.MsgType   `json:"type" binding:"required"`
	Mode        TopicCreateMode `json:"mode" binding:"required"`
}

type TopicRegisterResp struct {
	Topic     *model.Topic                `json:"topic"`
	Queues    []model.MessageQueue        `json:"queues"`
	Relations []model.BrokerQuereRelation `json:"relations"`
}
