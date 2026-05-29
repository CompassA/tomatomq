/*
 * @Author: Tomato
 * @Date: 2026-05-27 23:12:15
 * @LastEditTime: 2026-05-29 22:52:55
 */
package mqadmin

import (
	apperr "github.com/compassa/tomatomq/internal/admin/errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(req DatabaseRegisterReq) (*BrokerGroupDatabase, error) {
	database := buildNewDatabase(&req)

	id, err := s.repo.InsertOne(database)
	if err != nil {
		return nil, apperr.WrapError(apperr.DBErr, "db insert failed", err)
	}

	res, err := s.repo.QueryById(*id)
	if err != nil {
		return nil, apperr.WrapError(apperr.DBErr, "db select failed", err)
	}
	if len(res) == 0 {
		return nil, apperr.WrapError(apperr.DBErr, "query after insert failed", err)
	}

	return &res[0], nil
}

func (s *Service) QueryByBrokerGroup(brokerGroup string) ([]BrokerGroupDatabase, error) {
	res, err := s.repo.QueryByBrokerGroup(brokerGroup)
	if err != nil {
		return nil, apperr.WrapError(apperr.DBErr, "db select failed", err)
	}
	return res, nil
}
