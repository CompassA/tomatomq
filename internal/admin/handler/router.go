/*
 * @Author: Tomato
 * @Date: 2026-05-27 23:09:19
 * @LastEditTime: 2026-06-14 21:49:26
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
	// 反序列化
	var req mqadmin.DatabaseRegisterReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.Error(apperr.NewError(apperr.ParamInvalid, err.Error()))
		return
	}

	// 打印请求
	c.Set(midware.LogReqKey, req)

	// 业务逻辑
	model, err := h.adminsvc.Register(&req)
	if err != nil {
		c.Error(err)
		return
	}

	// 组装响应
	body := midware.OK(c, model)

	// 打印响应
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

func (h *Handler) TopicRegister(c *gin.Context) {
	var req mqadmin.TopicRegisterReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.Error(apperr.NewError(apperr.ParamInvalid, err.Error()))
		return
	}

	c.Set(midware.LogReqKey, req)

	res, err := h.adminsvc.CreateTopic(&req)
	if err != nil {
		c.Error(err)
		return
	}

	body := midware.OK(c, res)
	c.Set(midware.LogRespKey, body)
}
