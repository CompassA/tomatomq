/*
 * @Author: Tomato
 * @Date: 2026-05-27 23:16:31
 * @LastEditTime: 2026-05-28 00:34:03
 */
package handler

import (
	"net/http"

	"github.com/compassa/tomatomq/internal/admin/errors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool              `json:"success"`
	Data    any               `json:"data,omitempty"`
	Error   *errors.ErrorInfo `json:"error,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Fail(c *gin.Context, status int, err *errors.ErrorInfo) {
	c.JSON(status, Response{
		Success: false,
		Error:   err,
	})
}
