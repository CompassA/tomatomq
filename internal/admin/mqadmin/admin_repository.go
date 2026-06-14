/*
 * @Author: Tomato
 * @Date: 2026-05-21 21:37:35
 * @LastEditTime: 2026-06-15 00:14:58
 */

package mqadmin

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/compassa/tomatomq/internal/pkg/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type brokerAssignedNum struct {
	BrokerName string `gorm:"column:broker_name"`
	Cnt        int    `gorm:"column:cnt"`
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateTopicTx(
	topic *model.Topic,
	queue []model.MessageQueue,
	relation map[int]string,
) (resTopic *model.Topic, resQueue []model.MessageQueue, resRelation []model.BrokerQuereRelation, e error) {
	r.db.Transaction(func(tx *gorm.DB) error {
		resTopic, resQueue, resRelation, e = r.CreateTopic(tx, topic, queue, relation)
		return e
	})
	return
}

/**
 * 将内存中组装好的topic保存至DB
 * topic: 内存中组装好的topic, 无id
 * queue: 内存中组装好的队列, 无id、topicId
 * relation: [queue-index] -> [brokerName]
 */
func (r *Repository) CreateTopic(
	tx *gorm.DB,
	topic *model.Topic,
	queue []model.MessageQueue,
	relation map[int]string,
) (resTopic *model.Topic, resQueue []model.MessageQueue, resRelation []model.BrokerQuereRelation, e error) {
	// 创建topic
	resp := tx.Select("Name", "BrokerGroup", "Type", "QueueNum", "Status").Create(topic)
	if resp.Error != nil {
		e = fmt.Errorf("insert topic failed: %w", resp.Error)
		return
	}

	// 创建message_queue
	qmap := map[int]*model.MessageQueue{}
	for _, q := range queue {
		qmap[q.Index] = &q

		q.TopicId = topic.Id

		resp := tx.Select("TopicId", "Index", "DbId", "DBTableName", "Status").Create(&q)
		if resp.Error != nil {
			e = fmt.Errorf("insert queue failed: %w", resp.Error)
			return
		}
	}

	// 创建relation
	for index, brokerName := range relation {
		q, ok := qmap[index]
		if !ok {
			e = fmt.Errorf("queue not found in relation param, index=%d, brokerName=%s", index, brokerName)
			return
		}
		resp := tx.Select("BrokerGroup", "BrokerName", "QueueId").Create(&model.BrokerQuereRelation{
			BrokerGroup: topic.BrokerGroup,
			BrokerName:  brokerName,
			QueueId:     q.Id,
		})
		if resp.Error != nil {
			e = fmt.Errorf("insert relation failed: %w", resp.Error)
			return
		}
	}

	// 反查结果返回
	// topic
	tp, err := r.QueryTopicByNameTx(tx, topic.Name)
	if err != nil {
		e = fmt.Errorf("QueryTopicByNameTx: %w", err)
		return
	}
	if tp == nil {
		e = fmt.Errorf("QueryTopicByNameTx: topic %s not found", topic.Name)
		return
	}
	resTopic = tp

	// queue
	qs, err := r.QueryQueueByTopicIdTx(tx, tp.Id)
	if err != nil {
		e = fmt.Errorf("QueryQueueByTopicIdTx: %w", err)
		return
	}
	if len(qs) == 0 {
		e = fmt.Errorf("QueryQueueByTopicIdTx: queue not found, topicId %d", tp.Id)
		return
	}
	resQueue = qs

	// relation
	qids := make([]int64, 0, len(qs))
	for _, q := range qs {
		qids = append(qids, q.Id)
	}
	relations, err := r.QueryRelationByQueueIdTx(tx, qids)
	if err != nil {
		e = fmt.Errorf("QueryRelationByQueueIdTx: %w", err)
		return
	}
	if len(relations) == 0 {
		b, _ := json.Marshal(qids)
		e = fmt.Errorf("QueryRelationByQueueIdTx: queue not found, ququeIds %s", string(b))
		return
	}
	resRelation = relations
	return
}

func (r *Repository) InsertOne(database *model.BrokerGroupDatabase) (*int64, error) {
	return r.InsertOneTx(r.db, database)
}

func (r *Repository) InsertOneTx(db *gorm.DB, database *model.BrokerGroupDatabase) (*int64, error) {
	res := db.Select("Guid", "Dsn", "BrokerGroup").Create(database)
	if res.Error != nil {
		return nil, res.Error
	}
	return &database.Id, nil
}

func (r *Repository) QueryTopicByName(name string) (*model.Topic, error) {
	return r.QueryTopicByNameTx(r.db, name)
}

func (r *Repository) QueryTopicByNameTx(db *gorm.DB, name string) (*model.Topic, error) {
	res := []model.Topic{}
	resp := db.Raw(
		"SELECT id, name, broker_group, msg_type, msg_queue_num, status, gmt_create, gmt_modified FROM tomato_mq_topic WHERE name = @name",
		sql.Named("name", name)).Scan(&res)
	if resp.Error != nil {
		return nil, resp.Error
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (r *Repository) QueryQueueByTopicId(topicId int64) ([]model.MessageQueue, error) {
	return r.QueryQueueByTopicIdTx(r.db, topicId)
}

func (r *Repository) QueryQueueByTopicIdTx(db *gorm.DB, topicId int64) ([]model.MessageQueue, error) {
	res := []model.MessageQueue{}
	resp := db.Raw(
		"SELECT id, topic_id, `index`, db_id, table_name, status, gmt_create, gmt_modified FROM tomato_mq_mysql_queue WHERE topic_id = @topicId",
		sql.Named("topicId", topicId)).Scan(&res)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return res, nil
}

func (r *Repository) QueryRelationByQueueId(qIds []int64) ([]model.BrokerQuereRelation, error) {
	return r.QueryRelationByQueueIdTx(r.db, qIds)
}

func (r *Repository) QueryRelationByQueueIdTx(db *gorm.DB, qIds []int64) ([]model.BrokerQuereRelation, error) {
	if len(qIds) == 0 {
		return []model.BrokerQuereRelation{}, nil
	}

	res := []model.BrokerQuereRelation{}
	resp := db.Raw(
		"SELECT id, broker_group, broker_name, queue_id, gmt_create, gmt_modified FROM tomato_mq_broker_queue_relation WHERE queue_id IN(?)",
		qIds).Scan(&res)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return res, nil
}

func (r *Repository) QueryDBByBrokerGroup(brokerGroup string) ([]model.BrokerGroupDatabase, error) {
	return r.QueryDBByBrokerGroupTx(r.db, brokerGroup)
}

func (r *Repository) QueryDBByBrokerGroupTx(db *gorm.DB, brokerGroup string) ([]model.BrokerGroupDatabase, error) {
	res := []model.BrokerGroupDatabase{}

	resp := db.Raw(
		"SELECT id, db_guid, db_dsn, broker_group, gmt_create, gmt_modified FROM tomato_mq_db WHERE broker_group = @brokerGroup",
		sql.Named("brokerGroup", brokerGroup)).Scan(&res)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return res, nil
}

func (r *Repository) QueryBrokerAssignedNum(
	brokerGroup string,
	brokerNames map[string]struct{},
) ([]brokerAssignedNum, error) {
	return r.QueryBrokerAssignedNumTx(r.db, brokerGroup, brokerNames)
}

func (r *Repository) QueryBrokerAssignedNumTx(
	db *gorm.DB,
	brokerGroup string,
	brokerNames map[string]struct{},
) ([]brokerAssignedNum, error) {
	// DB分组查询, 统计每个BrokerName被分配了多少队列, 升序排序
	queryRes := []brokerAssignedNum{}
	resp := db.Raw(
		"SELECT broker_name, count(*) AS cnt FROM tomato_mq_broker_queue_relation WHERE broker_group = @brokerGroup GROUP BY broker_name ORDER BY cnt",
		sql.Named("brokerGroup", brokerGroup)).Scan(&queryRes)
	if resp.Error != nil {
		return nil, resp.Error
	}

	// 计算没被分配队列的broker
	zeroBrokerNames := map[string]struct{}{}
	for k := range brokerNames {
		zeroBrokerNames[k] = struct{}{}
	}
	for _, cnt := range queryRes {
		delete(zeroBrokerNames, cnt.BrokerName)
	}

	// 组装结果, 把没被分配队列的broker放到结果前面
	mergedRes := make([]brokerAssignedNum, 0, len(brokerNames))
	for zeroBrokerName := range zeroBrokerNames {
		mergedRes = append(mergedRes, brokerAssignedNum{
			BrokerName: zeroBrokerName,
			Cnt:        0,
		})
	}
	if len(queryRes) != 0 {
		mergedRes = append(mergedRes, queryRes...)
	}
	return mergedRes, nil
}

func (r *Repository) QueryDBById(id int64) ([]model.BrokerGroupDatabase, error) {
	return r.QueryDBByIdTx(r.db, id)
}

func (r *Repository) QueryDBByIdTx(db *gorm.DB, id int64) ([]model.BrokerGroupDatabase, error) {
	res := []model.BrokerGroupDatabase{}
	resp := db.Raw(
		"SELECT id, db_guid, db_dsn, broker_group, gmt_create, gmt_modified FROM tomato_mq_db WHERE id = @id",
		sql.Named("id", id)).Scan(&res)

	if resp.Error != nil {
		return nil, resp.Error
	}

	return res, nil
}
