/*
 * @Author: Tomato
 * @Date: 2026-05-27 23:09:19
 * @LastEditTime: 2026-05-28 23:34:28
 */
package handler

import (
	"net/http"

	apperr "github.com/compassa/tomatomq/internal/admin/errors"
	"github.com/compassa/tomatomq/internal/admin/mqadmin"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	adminsvc *mqadmin.Service
}

func NewHandler(adminsvc *mqadmin.Service) *Handler {
	return &Handler{
		adminsvc: adminsvc,
	}
}

func (h *Handler) DatabaseRegister(c *gin.Context) {
	var req mqadmin.DatabaseRegisterReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, apperr.NewError(apperr.PARAM_INVALID, err.Error()))
		return
	}

	model, err := h.adminsvc.Register(req)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err)
		return
	}

	OK(c, model)
}

func (h *Handler) DatabaseQueryByGroup(c *gin.Context) {
	brokerGroup := c.Query("group")
	if len(brokerGroup) == 0 {
		Fail(c, http.StatusBadRequest, apperr.NewError(apperr.PARAM_INVALID, "missing broker group"))
		return
	}

	res, err := h.adminsvc.QueryByBrokerGroup(brokerGroup)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err)
		return
	}

	OK(c, res)
}

func TopicRegister(c *gin.Context) {
}
