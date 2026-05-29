/*
 * @Author: Tomato
 * @Date: 2026-05-30 01:09:42
 * @LastEditTime: 2026-05-30 02:23:58
 */
package midware

import (
	"log/slog"
	"net/http"

	"github.com/compassa/tomatomq/internal/admin/config"
	apperr "github.com/compassa/tomatomq/internal/admin/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	TraceKey   = "_trace"
	LogReqKey  = "_log_req"
	LogRespKey = "_log_resq"
)

const (
	typeKey    = "type"
	traceIdKey = "traceId"
	dataKey    = "data"
	urlKey     = "url"
	errorKey   = "error"

	reqMark  = "req"
	respMark = "resp"

	headerTraceKey = "admin-http-trace-id"
)

func RequestUUIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		traceId := uuid.New().String()
		c.Set(TraceKey, traceId)
		c.Header(headerTraceKey, traceId)

		c.Next()
	}
}

func LogReqRespHandler() gin.HandlerFunc {
	logger := config.AppLogger
	return func(c *gin.Context) {
		c.Next()

		traceId := c.GetString(TraceKey)

		// 打印请求
		if req, ok := c.Get(LogReqKey); ok {
			logger.Info("http_log",
				slog.String(urlKey, c.Request.URL.Path),
				slog.String(typeKey, reqMark),
				slog.Any(dataKey, req),
				slog.String(traceIdKey, traceId))
		}

		// 打印响应
		if resp, ok := c.Get(LogRespKey); ok {
			logger.Info("http_log",
				slog.String(urlKey, c.Request.URL.Path),
				slog.String(typeKey, respMark),
				slog.Any(dataKey, resp),
				slog.String(traceIdKey, traceId))
		}
	}
}

func ErrorHandler() gin.HandlerFunc {
	logger := config.AppLogger
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}

		err := c.Errors.Last()
		if err == nil {
			return
		}

		traceId := c.GetString(TraceKey)

		// 打印错误
		logger.Error("http request error",
			slog.String(errorKey, err.Error()),
			slog.String(urlKey, c.Request.URL.Path),
			slog.String(traceIdKey, traceId))

		// 返回错误码
		errInfo, ok := err.Err.(*apperr.ErrorInfo)
		if ok {
			Fail(c, calcHttpStatus(errInfo.Code), errInfo)
		} else {
			Fail(c, http.StatusInternalServerError, apperr.NewError(apperr.UnknownServerErr, "internal server error"))
		}
	}
}

func calcHttpStatus(code apperr.ErrCode) int {
	if code < apperr.BizErr {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
