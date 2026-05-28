/*
 * @Author: Tomato
 * @Date: 2026-05-27 23:12:30
 * @LastEditTime: 2026-05-28 23:29:53
 */
package mqadmin

import (
	"fmt"
	"time"
)

type BrokerGroupDatabase struct {
	Id          int64     `json:"id" gorm:"column:id;primaryKey"`
	Guid        string    `json:"guid" gorm:"column:db_guid"`
	Dsn         string    `json:"dsn" gorm:"column:db_dsn"`
	BrokerGroup string    `json:"brokerGruop" gorm:"column:broker_group"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:gmt_create"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:gmt_modified"`
}

func (*BrokerGroupDatabase) TableName() string {
	return "tomato_mq_db"
}

type DatabaseRegisterReq struct {
	BrokerGroup string `json:"brokerGroup" binding:"required"`
	Name        string `json:"name" binding:"required"`
	User        string `json:"user" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Host        string `json:"host" binding:"required"`
	Port        int    `json:"port" binding:"gte=1025,lte=65535"`
}

func buildNewDatabase(req *DatabaseRegisterReq) *BrokerGroupDatabase {
	return &BrokerGroupDatabase{
		BrokerGroup: req.BrokerGroup,
		Dsn:         fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", req.User, req.Password, req.Host, req.Port, req.Name),
		Guid:        fmt.Sprintf("%s:%s", req.BrokerGroup, req.Name),
	}
}
