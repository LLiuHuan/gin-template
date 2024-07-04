// Package router
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:26
package router

import (
	"fmt"

	"github.com/LLiuHuan/gin-template/internal/alert"
	"github.com/LLiuHuan/gin-template/internal/metrics"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/cron"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/LLiuHuan/gin-template/internal/router/interceptor"
	"github.com/LLiuHuan/gin-template/pkg/errors"

	"go.uber.org/zap"
)

type resource struct {
	mux          core.Mux
	logger       *zap.Logger
	db           database.Repo
	cache        redis.Repo
	interceptors interceptor.Interceptor
	cronServer   cron.Server
	//ctx          context.Context
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

	//openBrowserUri := fmt.Sprintf("%s:%d", configs.Get().Project.Domain, configs.Get().Project.Port)

	//_, ok := file.IsExists(configs.ProjectInstallMark)
	//if !ok { // 未安装
	//	openBrowserUri += "/install"
	//} else { // 已安装
	// 初始化 DB
	dbRepo, err := database.New()
	if err != nil {
		logger.Fatal("new db err", zap.Error(err))
		panic(err)
	}
	r.db = dbRepo

	// 初始化 Cache
	cacheRepo, err := redis.New()
	if err != nil {
		logger.Fatal("new cache err", zap.Error(err))
		panic(err)
	}
	r.cache = cacheRepo

	// 初始化 CRON Server
	cronServer, err := cron.New(cronLogger, dbRepo, cacheRepo)
	if err != nil {
		logger.Fatal("new cron err", zap.Error(err))
		panic(err)
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
	r.interceptors = interceptor.New(logger, r.cache, r.db)

	// 设置 Render 路由
	//setRenderRouter(r)

	// 设置 API 路由
	setApiV1Router(r)

	// 设置 GraphQL 路由
	//setGraphQLRouter(r)

	// 设置 Socket 路由
	setSocketRouter(r)

	s := new(Server)
	s.Mux = mux
	s.Db = r.db
	s.Cache = r.cache
	s.CronServer = r.cronServer

	return s, nil
}

func ShutdownServer(logger *zap.Logger, s *Server) {
	if s.Db != nil {
		fmt.Println("db close")
		if err := s.Db.DBClose(); err != nil {
			logger.Error("db close err", zap.Error(err))
		}
	}

	if s.Cache != nil {
		fmt.Println("cache close")
		if err := s.Cache.Close(); err != nil {
			logger.Error("cache close err", zap.Error(err))
		}
	}

	if s.CronServer != nil {
		fmt.Println("cron close")
		s.CronServer.Stop()
	}
}
