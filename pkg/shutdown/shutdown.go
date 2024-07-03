// Package shutdown
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 22:15
package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

var _ Hook = (*hook)(nil)

// Hook 一个优雅的关闭钩子，默认带有 SIGINT 和 SIGTERM 信号
type Hook interface {
	// WithSignals 添加更多信号到钩子中
	WithSignals(signals ...syscall.Signal) Hook

	// Close 注册关闭句柄
	Close(funcs ...func())
}

type hook struct {
	ctx chan os.Signal
}

// NewHook 创建一个 Hook 实例
func NewHook() Hook {
	hook := &hook{
		ctx: make(chan os.Signal, 1),
	}

	return hook.WithSignals(syscall.SIGINT, syscall.SIGTERM)
}

// WithSignals 添加更多信号到钩子中
func (h *hook) WithSignals(signals ...syscall.Signal) Hook {
	for _, s := range signals {
		signal.Notify(h.ctx, s)
	}

	return h
}

// Close 注册关闭句柄
func (h *hook) Close(funcs ...func()) {
	select {
	case <-h.ctx:
	}
	signal.Stop(h.ctx)

	for _, f := range funcs {
		f()
	}
}
