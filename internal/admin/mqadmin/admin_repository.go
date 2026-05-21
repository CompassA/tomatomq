/*
 * @Author: Tomato
 * @Date: 2026-05-21 21:37:35
 * @LastEditTime: 2026-05-22 00:11:55
 */

package mqadmin

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type BrokerGroupDatabase struct {
	Id          int64     `gorm:"column:id"`
	Guid        string    `gorm:"column:db_guid"`
	Dsn         string    `gorm:"column:db_dsn"`
	BrokerGroup string    `gorm:"column:broker_group"`
	CreateAt    time.Time `gorm:"column:gmt_create"`
	UpdateAt    time.Time `gorm:"column:gmt_modified"`
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) QueryById(id int64) []BrokerGroupDatabase {
	res := []BrokerGroupDatabase{}
	r.db.Raw(`
		SELECT id, db_guid, db_dsn, broker_group, gmt_create, gmt_modified 
		FROM tomato_mq_db 
		WHERE id = @id`,
		sql.Named("id", id)).Scan(&res)

	return res
}
