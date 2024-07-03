//go:build windows

// Package kprocess
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-04 02:48
package kprocess

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var _ KProcess = (*kProcess)(nil)

type KProcess interface {
	i()

	Listen(network, addr string) (ln net.Listener, err error)
	Exit() <-chan struct{}
}

type kProcess struct {
	pidFile string
	pid     int
	ch      chan struct{}
	logger  *zap.Logger
}

func NewKProcess(logger *zap.Logger, pidFile string) KProcess {
	return &kProcess{
		logger:  logger,
		pidFile: pidFile,
	}
}

// Listen listens on the network address and returns an net.Listener.
// This shows how to use the upgrader
// with the graceful shutdown facilities of net/http.
// 这显示了如何使用升级程序
// 使用 net/http 的优雅关闭功能。
func (k *kProcess) Listen(network, addr string) (ln net.Listener, err error) {
	if k.ch == nil {
		k.ch = make(chan struct{})
	}
	k.pid = os.Getpid()
	k.logger.Info(fmt.Sprintf("exec process pid %d \n", k.pid))
	k.logger.Info("warning windows only support process shutdown ")

	go k.signal(k.stop)

	return net.Listen(network, addr)
}

func (k *kProcess) stop() error {
	close(k.ch)
	return nil
}

func (k *kProcess) upgrade() error {
	return nil
}

func (k *kProcess) Exit() <-chan struct{} {
	return k.ch
}

func (k *kProcess) signal(stopFunc func() error) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	for s := range sig {
		switch s {
		case syscall.SIGTERM:
			if stopFunc != nil {
				err := stopFunc()
				if err != nil {
					k.logger.Info(fmt.Sprintf("KProcess exec stopFunc failed:%v\n", err))
				}
				k.logger.Info(fmt.Sprintf("process [%d] stop...\n", k.pid))
			}
			return
		}
	}
}

func (k *kProcess) i() {}
