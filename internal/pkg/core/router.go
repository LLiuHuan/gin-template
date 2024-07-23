// Package core
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:29
package core

import (
	"time"

	_ "github.com/LLiuHuan/gin-template/docs"
	"github.com/LLiuHuan/gin-template/pkg/browser"
	"github.com/LLiuHuan/gin-template/pkg/env"
	"github.com/LLiuHuan/gin-template/pkg/errors"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// 封装handler，统一处理
func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		//handler := handler
		funcs[i] = func(c *gin.Context) {
			ctx := newContext(c)
			defer releaseContext(ctx)

			handler(ctx)
		}
	}

	return funcs
}

// NewRouter 创建路由
func NewRouter(logger *zap.Logger, options ...Option) (Mux, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	muxEngine := &mux{
		engine: gin.New(),
	}

	// 静态文件
	//mux.engine.StaticFS("assets", http.FS(assets.Bootstrap))
	//mux.engine.SetHTMLTemplate(template.Must(template.New("").ParseFS(assets.Templates, "templates/**/*")))

	opt := new(option)
	for _, f := range options {
		f(opt)
	}

	if !opt.disablePProf {
		// pprof
		if !env.Active().IsPro() {
			pprof.Register(muxEngine.engine) // register pprof to gin
		}
	}

	if !opt.disableSwagger {
		// swagger
		if !env.Active().IsPro() {
			muxEngine.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // register swagger
		}
	}

	if !opt.disablePrometheus {
		// prometheus
		muxEngine.engine.GET("/debug/metrics", gin.WrapH(promhttp.Handler())) // register prometheus
	}

	if opt.enableCors {
		// 跨域
		muxEngine.engine.Use(MiddlewareCors())
	}

	if opt.enableOpenBrowser != "" {
		_ = browser.Open(opt.enableOpenBrowser)
	}

	// recover两次，防止处理时发生panic，尤其是在OnPanicNotify中。
	muxEngine.engine.Use(MiddlewareRecover(logger))

	muxEngine.engine.Use(MiddlewareTrace(logger, opt))

	// 限流
	if opt.enableRate {
		muxEngine.engine.Use(MiddlewareLimit())
	}

	//
	muxEngine.engine.NoMethod(wrapHandlers(DisableTraceLog)...)
	muxEngine.engine.NoRoute(wrapHandlers(DisableTraceLog)...)

	system := muxEngine.Group("/system")
	{
		// 健康检查
		system.GET("/health", func(ctx Context) {
			resp := &struct {
				Timestamp   time.Time `json:"timestamp"`
				Environment string    `json:"environment"`
				Host        string    `json:"host"`
				Status      string    `json:"status"`
			}{
				Timestamp:   time.Now(),
				Environment: env.Active().Value(),
				Host:        ctx.Host(),
				Status:      "ok",
			}
			ctx.Payload(resp)
		})
	}

	return muxEngine, nil
}
