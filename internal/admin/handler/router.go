/*
 * @Author: Tomato
 * @Date: 2026-05-27 23:09:19
 * @LastEditTime: 2026-05-30 02:21:34
 */
package handler

import (
	apperr "github.com/compassa/tomatomq/internal/admin/errors"
	"github.com/compassa/tomatomq/internal/admin/midware"
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
		c.Error(apperr.NewError(apperr.ParamInvalid, err.Error()))
		return
	}

	c.Set(midware.LogReqKey, req)
	model, err := h.adminsvc.Register(req)
	if err != nil {
		c.Error(err)
		return
	}

	body := midware.OK(c, model)

	c.Set(midware.LogRespKey, body)
}

func (h *Handler) DatabaseQueryByGroup(c *gin.Context) {
	brokerGroup := c.Query("group")
	if len(brokerGroup) == 0 {
		c.Error(apperr.NewError(apperr.ParamInvalid, "missing broker group"))
		return
	}
	c.Set(midware.LogReqKey, brokerGroup)

	res, err := h.adminsvc.QueryByBrokerGroup(brokerGroup)
	if err != nil {
		c.Error(err)
		return
	}

	body := midware.OK(c, res)

	c.Set(midware.LogRespKey, body)
}

func TopicRegister(c *gin.Context) {
}
