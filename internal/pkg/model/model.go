/*
 * @Author: Tomato
 * @Date: 2026-06-14 17:18:14
 * @LastEditTime: 2026-06-14 22:12:30
 */
package model

import "time"

/**
 * DB领域模型
 */
type BrokerGroupDatabase struct {
	Id          int64     `json:"id" gorm:"column:id;primaryKey"`         // DBId
	Guid        string    `json:"guid" gorm:"column:db_guid"`             // 有业务含义的唯一ID, 拼接格式: "${BrokerGroup}:${数据库名称}"
	Dsn         string    `json:"dsn" gorm:"column:db_dsn"`               // 数据库资源链接
	BrokerGroup string    `json:"brokerGruop" gorm:"column:broker_group"` // 资源分组
	CreatedAt   time.Time `json:"createdAt" gorm:"column:gmt_create"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:gmt_modified"`
}

func (*BrokerGroupDatabase) TableName() string {
	return "tomato_mq_db"
}

/**
 * Topic领域模型
 */
type Topic struct {
	Id          int64     `json:"id" gorm:"column:id;primaryKey"`         // TopicId
	Name        string    `json:"name" gorm:"column:name"`                // Topic名称
	BrokerGroup string    `json:"brokerGroup" gorm:"column:broker_group"` // Topic资源分组
	Type        MsgType   `json:"type" gorm:"column:msg_type"`            // Topic的消息类型
	QueueNum    int       `json:"queueNum" gorm:"column:msg_queue_num"`   // Topic有多少个分区, 分区数量由BrokerGroup持有的DB库数量决定, 有几个库就会创建几个分区
	Status      string    `json:"status" gorm:"column:status"`            // Topic创建状态, 预留字段
	CreatedAt   time.Time `json:"createdAt" gorm:"column:gmt_create"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:gmt_modified"`
}

func (*Topic) TableName() string {
	return "tomato_mq_topic"
}

// topic消息类型
type MsgType int

const (
	UnknownMsgType MsgType = iota
	Normal                 // 普通消息
	Ordered                // 顺序消息
	Trx                    // 事务消息
	Scheduled              // 定时消息
)

/**
 * 消息队列领域模型
 */
type MessageQueue struct {
	Id          int64     `json:"id" gorm:"column:id"`                // 队列Id
	TopicId     int64     `json:"topicId" gorm:"column:topic_id"`     // 所属的Topic
	Index       int       `json:"index" gorm:"column:index"`          // 队列数字编号 0-"N-1"
	DbId        int64     `json:"dbID" gorm:"column:db_id"`           // 消息保存在哪个DB中
	DBTableName string    `json:"tableName" gorm:"column:table_name"` // 消息保存在DB的哪张表中
	Status      string    `json:"status" gorm:"column:status"`        // 队列创建状态, 预留字段
	CreatedAt   time.Time `json:"createdAt" gorm:"column:gmt_create"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:gmt_modified"`
}

func (*MessageQueue) TableName() string {
	return "tomato_mq_mysql_queue"
}

/**
 * 队列与Mq的分配关系
 */
type BrokerQuereRelation struct {
	Id          int64     `json:"id" gorm:"column:id"`                    // id
	BrokerGroup string    `json:"brokerGroup" gorm:"column:broker_group"` // 所属的资源分组
	BrokerName  string    `json:"brokerName" gorm:"column:broker_name"`   // Broker名称
	QueueId     int64     `json:"queueId" gorm:"column:queue_id"`         // Broker管理的队列
	CreatedAt   time.Time `json:"createdAt" gorm:"column:gmt_create"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:gmt_modified"`
}

func (*BrokerQuereRelation) TableName() string {
	return "tomato_mq_broker_queue_relation"
}
