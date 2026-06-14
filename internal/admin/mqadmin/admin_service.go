/*
 * @Author: Tomato
 * @Date: 2026-05-27 23:12:15
 * @LastEditTime: 2026-06-14 23:32:51
 */
package mqadmin

import (
	"fmt"

	apperr "github.com/compassa/tomatomq/internal/admin/errors"
	"github.com/compassa/tomatomq/internal/admin/meta"
	"github.com/compassa/tomatomq/internal/pkg/model"
)

type Service struct {
	repo     *Repository
	etcdRepo *meta.BrokerCacheRepo
}

func NewService(repo *Repository, etcdRepo *meta.BrokerCacheRepo) *Service {
	return &Service{
		repo:     repo,
		etcdRepo: etcdRepo,
	}
}

func (s *Service) Register(req *DatabaseRegisterReq) (*model.BrokerGroupDatabase, error) {
	database := buildNewDatabase(req)

	id, err := s.repo.InsertOne(database)
	if err != nil {
		return nil, apperr.WrapError(apperr.DBErr, "db insert failed", err)
	}

	res, err := s.repo.QueryDBById(*id)
	if err != nil {
		return nil, apperr.WrapError(apperr.DBErr, "db select failed", err)
	}
	if len(res) == 0 {
		return nil, apperr.WrapError(apperr.DBErr, "query after insert failed", err)
	}

	return &res[0], nil
}

func (s *Service) CreateTopic(req *TopicRegisterReq) (*TopicRegisterResp, error) {
	// 读取BrokerGroup下的db资源
	dbs, err := s.repo.QueryDBByBrokerGroup(req.BrokerGroup)
	if err != nil {
		return nil, apperr.WrapError(apperr.DBErr, "CreateTopic->QueryDBByBrokerGroup", err)
	}

	// 读取BrokerGroup下活跃的broker
	brokers := s.etcdRepo.GetBrokerByGroup(req.BrokerGroup)
	if len(brokers) == 0 {
		return nil, apperr.NewError(apperr.BizErr, fmt.Sprintf("no active broker in group[%s]", req.BrokerGroup))
	}

	// 读取每个Broker已经分配的队列数, 按队列数升序排序
	brokerNameMap := make(map[string]struct{}, len(brokers))
	for _, broker := range brokers {
		brokerNameMap[broker.Name] = struct{}{}
	}
	brokerAssigned, err := s.repo.QueryBrokerAssignedNum(req.BrokerGroup, brokerNameMap)
	if err != nil {
		return nil, apperr.WrapError(apperr.DBErr, "CreateTopic->QueryBrokerAssignedNum", err)
	}

	// Normal模式, queue平均分配给Broker; Restrict模式, broker数量必须大于queue数量, 每个queue分配不同的broker
	brokerNum := len(brokerNameMap)
	queueNum := len(dbs)
	if req.Mode == Restrict && brokerNum < queueNum {
		return nil, apperr.NewError(apperr.BizErr, fmt.Sprintf("topic create failed in restrict mode, broker num %d less than db num %d", brokerNum, queueNum))
	}
	queue := make([]model.MessageQueue, 0, queueNum)
	relation := make(map[int]string, queueNum)
	for i, db := range dbs {
		queue = append(queue, model.MessageQueue{
			Index:       i,
			DbId:        db.Id,
			DBTableName: fmt.Sprintf("%s_%s_%d", req.BrokerGroup, req.TopicName, i),
			Status:      "",
		})

		relation[i] = brokerAssigned[i%brokerNum].BrokerName
	}

	// 保存
	tp, q, r, err := s.repo.CreateTopicTx(&model.Topic{
		Name:        req.TopicName,
		BrokerGroup: req.BrokerGroup,
		Type:        req.Type,
		QueueNum:    queueNum,
		Status:      "",
	}, queue, relation)
	if err != nil {
		return nil, apperr.WrapError(apperr.BizErr, "save topic failed", err)
	}

	return &TopicRegisterResp{
		Topic:     tp,
		Queues:    q,
		Relations: r,
	}, nil
}

func (s *Service) QueryByBrokerGroup(brokerGroup string) ([]model.BrokerGroupDatabase, error) {
	res, err := s.repo.QueryDBByBrokerGroup(brokerGroup)
	if err != nil {
		return nil, apperr.WrapError(apperr.DBErr, "db select failed", err)
	}
	return res, nil
}
