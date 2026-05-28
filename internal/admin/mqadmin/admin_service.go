/*
 * @Author: Tomato
 * @Date: 2026-05-27 23:12:15
 * @LastEditTime: 2026-05-28 23:19:21
 */
package mqadmin

import (
	"log/slog"

	"github.com/compassa/tomatomq/internal/admin/config"
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

func (s *Service) Register(req DatabaseRegisterReq) (*BrokerGroupDatabase, *apperr.ErrorInfo) {
	database := buildNewDatabase(&req)

	id, err := s.repo.InsertOne(database)
	if err != nil {
		config.AppLogger.Error("database insert failed", slog.String("error", err.Error()))
		return nil, apperr.NewError(apperr.DB_ERR, "db insert failed")
	}

	res, err := s.repo.QueryById(*id)
	if err != nil {
		config.AppLogger.Error("database select failed", slog.String("error", err.Error()))
		return nil, apperr.NewError(apperr.DB_ERR, "db select failed")
	}
	if len(res) == 0 {
		return nil, apperr.NewError(apperr.DB_ERR, "query after insert failed")
	}

	return &res[0], nil
}

func (s *Service) QueryByBrokerGroup(brokerGroup string) ([]BrokerGroupDatabase, *apperr.ErrorInfo) {
	res, err := s.repo.QueryByBrokerGroup(brokerGroup)
	if err != nil {
		config.AppLogger.Error("database select failed", slog.String("error", err.Error()))
		return nil, apperr.NewError(apperr.DB_ERR, "db select failed")
	}
	return res, nil
}
