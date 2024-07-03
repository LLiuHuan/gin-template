//go:build darwin

// Package kprocess
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-04 01:34
package kprocess

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"go.uber.org/zap"
)

var _ KProcess = (*kProcess)(nil)

type KProcess interface {
	i()

	Listen(network, addr string) (ln net.Listener, err error)
	Exit() <-chan struct{}
}

type kProcess struct {
	pidFile   string
	pid       int
	processUp *tableflip.Upgrader
	logger    *zap.Logger
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
	k.pid = os.Getpid()
	k.logger.Info(fmt.Sprintf("exec process pid %d \n", k.pid))

	k.processUp, err = tableflip.New(tableflip.Options{
		UpgradeTimeout: 5 * time.Second,
		PIDFile:        k.pidFile,
	})
	if err != nil {
		return nil, err
	}

	go k.signal(k.upgrade, k.stop)

	// Listen 必须在 Ready 之前调用
	// Listen must be called before Ready
	if network != "" && addr != "" {
		ln, err = k.processUp.Listen(network, addr)
		if err != nil {
			return nil, err
		}
	}
	if err := k.processUp.Ready(); err != nil {
		return nil, err
	}

	return ln, nil
}

func (k *kProcess) stop() error {
	if k.processUp != nil {
		k.processUp.Stop()
		return os.Remove(k.pidFile)
	}
	return nil
}

func (k *kProcess) upgrade() error {
	if k.processUp != nil {
		return k.processUp.Upgrade()
	}
	return nil
}

func (k *kProcess) Exit() <-chan struct{} {
	if k.processUp != nil {
		return k.processUp.Exit()
	}
	ch := make(chan struct{})
	close(ch)
	return ch
}

func (k *kProcess) signal(upgradeFunc, stopFunc func() error) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM, os.Kill, os.Interrupt)
	for s := range sig {
		switch s {
		case syscall.SIGTERM, os.Kill, os.Interrupt:
			if stopFunc != nil {
				err := stopFunc()
				if err != nil {
					k.logger.Info(fmt.Sprintf("KProcess exec stopFunc failed:%v\n", err))
				}
				k.logger.Info(fmt.Sprintf("process [%d] stop...\n", k.pid))
			}
			return
		case syscall.SIGUSR1, syscall.SIGUSR2:
			if upgradeFunc != nil {
				err := upgradeFunc()
				if err != nil {
					k.logger.Info(fmt.Sprintf("KProcess exec Upgrade failed:%v\n", err))
				}
				k.logger.Info(fmt.Sprintf("process [%d] restart...\n", k.pid))
			}
		}
	}
}

func (k *kProcess) i() {}
