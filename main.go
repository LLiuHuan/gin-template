// Package gin_template
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 18:07
package main

import (
	"context"
	"fmt"
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/router"
	"github.com/LLiuHuan/gin-template/pkg/app"
	"github.com/LLiuHuan/gin-template/pkg/env"
	"github.com/LLiuHuan/gin-template/pkg/grace"
	"github.com/LLiuHuan/gin-template/pkg/kprocess"
	"github.com/LLiuHuan/gin-template/pkg/logger"
	"github.com/LLiuHuan/gin-template/pkg/timeutil"
)

// 初始化执行
func init() {
	//gin.SetMode(gin.ReleaseMode)

}

// @title swagger 接口文档
// @version 2.0
// @description

// @contact.name
// @contact.url
// @contact.email

// @license.name MIT
// @license.url https://github.com/LLiuHuan/gin-template/blob/master/LICENSE

// @securityDefinitions.apikey  LoginToken
// @in                          header
// @name                        token
func main() {
	accessLogger, err := logger.NewJSONLogger(
		logger.WithDisableConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileRotationP(configs.ProjectLogFile))
	if err != nil {
		panic(err)
	}

	// 初始化 cron logger
	cronLogger, err := logger.NewJSONLogger(
		logger.WithDisableConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileP(configs.ProjectCronLogFile),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = accessLogger.Sync()
		_ = cronLogger.Sync()
	}()

	// 初始化 HTTP 服务
	s, err := router.NewHTTPServer(accessLogger, cronLogger)
	if err != nil {
		panic(err)
	}

	server := grace.NewServer(fmt.Sprintf("%s:%d", configs.Get().Project.Domain, configs.Get().Project.Port), s.Mux, s.Opts...)
	kp := kprocess.NewKProcess(accessLogger, configs.Get().Project.PidFile)
	ln, err := kp.Listen("tcp", server.Addr)
	if err != nil {
		panic(err)
	}

	serverCloe := make(chan struct{})
	go func() {
		defer func() {
			close(serverCloe)
		}()
		err = server.ServeWithListener(ln)
		if err != nil {
			accessLogger.Info(fmt.Sprintf("App run Serve: %v\n", err))
		}
	}()

	select {
	case <-kp.Exit():
	case <-serverCloe:
	}

	// Make sure to set a deadline on exiting the process
	// after upg.Exit() is closed. No new upgrades can be
	// performed if the parent doesn't exit.
	app.AppPrepareForceExit(accessLogger)

	err = server.Shutdown(context.Background())
	if err != nil {
		accessLogger.Info(fmt.Sprintf("App run Shutdown: %v\n", err))
	}

	accessLogger.Info("App server Shutdown ok")
}
