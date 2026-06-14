/*
 * @Author: Tomato
 * @Date: 2026-05-28 00:33:06
 * @LastEditTime: 2026-06-14 22:25:41
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
	return nil
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
