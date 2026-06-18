/*
 * @Author: Tomato
 * @Date: 2026-05-28 00:33:06
 * @LastEditTime: 2026-06-19 01:48:26
 */
package errors

import (
	"fmt"
)

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
	errstr := ""
	if e.Err != nil {
		errstr = e.Err.Error()
	}
	return fmt.Sprintf("Code:%d,Message:%s,InnerErr:%s", e.Code, e.Message, errstr)
}

func (e *ErrorInfo) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return nil
}

func NewError(code ErrCode, msg string) *ErrorInfo {
	return WrapError(code, msg, nil)
}

func WrapError(code ErrCode, msg string, err error) *ErrorInfo {
	return &ErrorInfo{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}
