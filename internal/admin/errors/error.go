/*
 * @Author: Tomato
 * @Date: 2026-05-28 00:33:06
 * @LastEditTime: 2026-05-30 02:25:21
 */
package errors

import "strconv"

type ErrCode int

const (
	ParamInvalid     ErrCode = 10001
	BizErr           ErrCode = 20001
	DBErr            ErrCode = 30002
	UnknownServerErr ErrCode = 100000
)

type ErrorInfo struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
	Err     error
}

func (e *ErrorInfo) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "Code:" + strconv.Itoa(int(e.Code)) + ",Message:" + e.Message
}

func (e *ErrorInfo) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return e
}

func NewError(code ErrCode, msg string) *ErrorInfo {
	return &ErrorInfo{
		Code:    code,
		Message: msg,
	}
}

func WrapError(code ErrCode, msg string, err error) *ErrorInfo {
	return &ErrorInfo{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}
