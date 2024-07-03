// Package errors
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:36
package errors

import (
	"fmt"
	"io"
	"runtime"

	"github.com/pkg/errors"
)

func callers() []uintptr {
	var pcs [32]uintptr
	l := runtime.Callers(3, pcs[:])
	return pcs[:l]
}

// Error 调用者堆栈信息错误
type Error interface {
	error
	t()
}

var _ Error = (*item)(nil)
var _ fmt.Formatter = (*item)(nil)

type item struct {
	msg   string
	stack []uintptr
}

func (i *item) Error() string {
	return i.msg
}

func (i *item) t() {}

// Format 由 go.uber.org/zap 以 Verbose 形式使用
func (i *item) Format(s fmt.State, verb rune) {
	io.WriteString(s, i.msg)
	io.WriteString(s, "\n")

	for _, pc := range i.stack {
		fmt.Fprintf(s, "%+v\n", errors.Frame(pc))
	}
}

// New 创建一个新错误
func New(msg string) Error {
	return &item{msg: msg, stack: callers()}
}

// Errorf 创建一个新错误
func Errorf(format string, args ...interface{}) Error {
	return &item{msg: fmt.Sprintf(format, args...), stack: callers()}
}

// Wrap 错误中包含一些额外的消息
func Wrap(err error, msg string) Error {
	if err == nil {
		return nil
	}

	e, ok := err.(*item)
	if !ok {
		return &item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), stack: callers()}
	}

	e.msg = fmt.Sprintf("%s; %s", msg, e.msg)
	return e
}

// Wrapf 错误中包含一些额外的消息
func Wrapf(err error, format string, args ...interface{}) Error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args...)

	e, ok := err.(*item)
	if !ok {
		return &item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), stack: callers()}
	}

	e.msg = fmt.Sprintf("%s; %s", msg, e.msg)
	return e
}

// WithStack 添加调用者堆栈信息
func WithStack(err error) Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*item); ok {
		return e
	}

	return &item{msg: err.Error(), stack: callers()}
}
