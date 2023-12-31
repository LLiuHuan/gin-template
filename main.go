// Package gin_template
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-08-16 11:03
package main

import (
	"errors"
	"fmt"
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/router"
	"github.com/LLiuHuan/gin-template/pkg/env"
	"github.com/LLiuHuan/gin-template/pkg/grace"
	"github.com/LLiuHuan/gin-template/pkg/logger"
	"github.com/LLiuHuan/gin-template/pkg/timeutil"
	"go.uber.org/zap"
	"net/http"
)

// 初始化执行
func init() {

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
	accessLogger, err := logger.NewLogger(
		logger.WithDisableConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileRotationP(configs.ProjectLogFile))
	if err != nil {
		panic(err)
	}

	// 初始化 cron logger
	cronLogger, err := logger.NewLogger(
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

	s, err := router.NewHTTPServer(accessLogger, cronLogger)
	if err != nil {
		panic(err)
	}
	server := grace.NewServer(fmt.Sprintf("%s:%d", configs.Get().Project.Domain, configs.Get().Project.Port), s.Mux, s.Opts...)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		accessLogger.Fatal("http server startup err", zap.Error(err))
	}
}
