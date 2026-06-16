/*
 * @Author: Tomato
 * @Date: 2026-05-30 01:09:42
 * @LastEditTime: 2026-06-17 00:31:45
 */
package midware

import (
	"context"
	"log/slog"
	"net/http"

	apperr "github.com/compassa/tomatomq/internal/admin/errors"
	"github.com/compassa/tomatomq/pkg/tomatolog"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 参数key
const (
	// gin request context 参数key
	CtxTraceKey = "__trace__" // 本次请求的trace_id

	// gin ctx参数
	LogReqKey  = "_log_req"
	LogRespKey = "_log_resq"

	// 日志全局参数
	traceKey = "traceId"

	// 日志参数
	typeKey  = "type"  // 区分请求日志和响应日志. req 请求, resp 响应
	dataKey  = "data"  // 请求数据\响应数据
	urlKey   = "url"   // 请求的url
	errorKey = "error" // 错误码

	// 数据
	reqMark  = "req"  // typeKey的值, 代表请求
	respMark = "resp" // typeKey的值, 代表响应

	// http header key
	headerTraceKey = "X-trace-id" // http请求中的traceIDkey
)

func RequestUUIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		traceId := c.GetHeader(headerTraceKey)
		if len(traceId) == 0 {
			traceId = uuid.New().String()
		}
		c.Header(headerTraceKey, traceId)

		// 绑定日志参数
		logger := slog.Default().With(traceKey, traceId)

		// 绑定logger、traceId
		ctx := context.WithValue(c.Request.Context(), CtxTraceKey, traceId)
		ctx = context.WithValue(ctx, tomatolog.CtxLoggerKey, logger)
		c.Request = c.Request.WithContext(ctx)

		// 执行逻辑
		c.Next()
	}
}

func LogReqRespHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行业务逻辑
		c.Next()

		// 执行完业务逻辑后, 打印请求与响应
		logger := tomatolog.FromCtx(c.Request.Context())

		if req, ok := c.Get(LogReqKey); ok {
			logger.Info("http request",
				slog.String(urlKey, c.Request.URL.Path),
				slog.String(typeKey, reqMark),
				slog.Any(dataKey, req))
		}

		// 打印响应
		if resp, ok := c.Get(LogRespKey); ok {
			logger.Info("http response",
				slog.String(urlKey, c.Request.URL.Path),
				slog.String(typeKey, respMark),
				slog.Any(dataKey, resp))
		}
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 执行业务逻辑
		c.Next()

		if c.Writer.Written() {
			return
		}

		err := c.Errors.Last()
		if err == nil {
			return
		}

		// 打印错误
		logger := tomatolog.FromCtx(c.Request.Context())

		logger.Error("http request error",
			slog.String(errorKey, err.Error()),
			slog.String(urlKey, c.Request.URL.Path))

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
