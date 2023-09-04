// Package router
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-08-16 16:00
// @description: 路由
package router

import (
	"errors"
	"github.com/LLiuHuan/gin-template/internal/alert"
	"github.com/LLiuHuan/gin-template/internal/metrics"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/cron"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/LLiuHuan/gin-template/internal/router/interceptor"
	"go.uber.org/zap"
)

type resource struct {
	mux          core.Mux
	logger       *zap.Logger
	db           database.Repo
	cache        redis.Repo
	interceptors interceptor.Interceptor
	cronServer   cron.Server
}

type Server struct {
	Mux        core.Mux
	Db         database.Repo
	Cache      redis.Repo
	CronServer cron.Server
}

func NewHTTPServer(logger *zap.Logger, cronLogger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	r := new(resource)
	r.logger = logger

	//openBrowserUri := fmt.Sprintf("http://%s:%d", configs.Get().Project.Domain, configs.Get().Project.Port)

	// TODO: 后续判断一下是否安装，如果没安装可以提示让跳转到初始化页面
	//_, ok := file.IsExists(configs.ProjectInstallMark)
	//if !ok { // 未安装
	//	openBrowserUri += "/install"
	//} else { // 已安装
	// 初始化 DB
	dbRepo, err := database.New()
	if err != nil {
		logger.Fatal("new db err", zap.Error(err))
	}
	r.db = dbRepo

	// 初始化 Cache
	cacheRepo, err := redis.New()
	if err != nil {
		logger.Fatal("new cache err", zap.Error(err))
	}
	r.cache = cacheRepo

	// 初始化 CRON Server
	cronServer, err := cron.New(cronLogger, dbRepo, cacheRepo)
	if err != nil {
		logger.Fatal("new cron err", zap.Error(err))
	}
	cronServer.Start()
	r.cronServer = cronServer
	//}

	mux, err := core.NewRouter(logger,
		//core.WithEnableOpenBrowser(openBrowserUri),
		core.WithEnableCors(),
		core.WithEnableRate(),
		core.WithAlertNotify(alert.NotifyHandler(logger)),
		core.WithRecordMetrics(metrics.RecordHandler(logger)),
	)

	if err != nil {
		panic(err)
	}

	r.mux = mux
	// 拦截器
	r.interceptors = interceptor.New(logger, r.cache, r.db)

	setRenderRouter(r)

	setApiRouter(r)

	s := new(Server)
	s.Mux = mux
	s.Db = r.db
	s.Cache = r.cache
	s.CronServer = r.cronServer

	return s, nil
}
