/*
 * @Author: Tomato
 * @Date: 2026-05-28 00:33:06
 * @LastEditTime: 2026-05-28 00:36:11
 */
package errors

type ErrCode int

const (
	PARAM_INVALID ErrCode = 10001
	BIZ_ERR       ErrCode = 20001
	DB_ERR        ErrCode = 20002
)

type ErrorInfo struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func NewError(code ErrCode, msg string) *ErrorInfo {
	return &ErrorInfo{
		Code:    code,
		Message: msg,
	}
}
