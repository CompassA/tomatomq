/*
 * @Author: Tomato
 * @Date: 2026-05-21 21:37:35
 * @LastEditTime: 2026-05-28 23:17:35
 */

package mqadmin

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) InsertOne(database *BrokerGroupDatabase) (*int64, error) {
	res := r.db.Select("Guid", "Dsn", "BrokerGroup").Create(database)
	if res.Error != nil {
		return nil, res.Error
	}
	return &database.Id, nil
}

func (r *Repository) QueryByBrokerGroup(brokerGroup string) ([]BrokerGroupDatabase, error) {
	res := []BrokerGroupDatabase{}

	resp := r.db.Raw(`
		SELECT id, db_guid, db_dsn, broker_group, gmt_create, gmt_modified 
		FROM tomato_mq_db 
		WHERE broker_group = @brokerGroup`,
		sql.Named("brokerGroup", brokerGroup)).Scan(&res)
	if resp.Error != nil {
		return nil, fmt.Errorf("query failed: %w", resp.Error)
	}

	return res, nil
}

func (r *Repository) QueryById(id int64) ([]BrokerGroupDatabase, error) {
	res := []BrokerGroupDatabase{}
	resp := r.db.Raw(`
		SELECT id, db_guid, db_dsn, broker_group, gmt_create, gmt_modified 
		FROM tomato_mq_db 
		WHERE id = @id`,
		sql.Named("id", id)).Scan(&res)

	if resp.Error != nil {
		return nil, fmt.Errorf("query failed: %w", resp.Error)
	}

	return res, nil
}
