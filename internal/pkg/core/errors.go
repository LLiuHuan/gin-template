// Package core
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:35
package core

import (
	"github.com/LLiuHuan/gin-template/pkg/env"
	"github.com/LLiuHuan/gin-template/pkg/errors"
)

var _ BusinessError = (*businessError)(nil)

type BusinessError interface {
	// i 为了避免被其他包实现
	i()

	// WithError 设置错误信息
	WithError(err error) BusinessError

	// WithAlert 设置告警通知
	WithAlert() BusinessError

	// BusinessCode 获取业务码
	BusinessCode() int

	// HTTPCode 获取 HTTP 状态码
	HTTPCode() int

	// Message 获取错误描述
	Message() string

	// StackError 获取带堆栈的错误信息
	StackError() error

	// IsAlert 是否开启告警通知
	IsAlert() bool
}

type businessError struct {
	httpCode     int    // HTTP 状态码
	businessCode int    // 业务码
	message      string // 错误描述
	stackError   error  // 含有堆栈信息的错误
	isAlert      bool   // 是否告警通知
}

func Error(httpCode, businessCode int, message string) BusinessError {
	return &businessError{
		httpCode:     httpCode,
		businessCode: businessCode,
		message:      message,
		isAlert:      false,
	}
}

func (e *businessError) i() {}

func (e *businessError) WithError(err error) BusinessError {
	e.stackError = errors.WithStack(err)
	return e
}

func (e *businessError) WithAlert() BusinessError {
	e.isAlert = true
	return e
}

func (e *businessError) HTTPCode() int {
	return e.httpCode
}

func (e *businessError) BusinessCode() int {
	return e.businessCode
}

func (e *businessError) Message() string {
	if env.Active().IsDev() {
		return e.message + ":" + e.stackError.Error()
	}
	return e.message
}

func (e *businessError) StackError() error {
	return e.stackError
}

func (e *businessError) IsAlert() bool {
	return e.isAlert
}
