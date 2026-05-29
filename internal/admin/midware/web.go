/*
 * @Author: Tomato
 * @Date: 2026-05-27 23:16:31
 * @LastEditTime: 2026-05-30 01:54:49
 */
package midware

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

func OK(c *gin.Context, data any) *Response {
	body := Response{
		Success: true,
		Data:    data,
	}
	c.JSON(http.StatusOK, body)
	return &body
}

func Fail(c *gin.Context, status int, err *errors.ErrorInfo) *Response {
	body := Response{
		Success: false,
		Error:   err,
	}
	c.JSON(status, body)
	return &body
}
