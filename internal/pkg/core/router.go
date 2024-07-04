// Package core
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:29
package core

import (
	"fmt"
	"time"

	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/pkg/browser"
	"github.com/LLiuHuan/gin-template/pkg/color"
	"github.com/LLiuHuan/gin-template/pkg/env"
	"github.com/LLiuHuan/gin-template/pkg/errors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// see https://patorjk.com/software/taag/#p=testall&f=Graffiti&t=gin-template
var _UI = fmt.Sprintf(`
     ██████╗ ██╗███╗   ██╗   ████████╗███████╗███╗   ███╗██████╗ ██╗      █████╗ ████████╗███████╗
    ██╔════╝ ██║████╗  ██║   ╚══██╔══╝██╔════╝████╗ ████║██╔══██╗██║     ██╔══██╗╚══██╔══╝██╔════╝
    ██║  ███╗██║██╔██╗ ██║█████╗██║   █████╗  ██╔████╔██║██████╔╝██║     ███████║   ██║   █████╗  
    ██║   ██║██║██║╚██╗██║╚════╝██║   ██╔══╝  ██║╚██╔╝██║██╔═══╝ ██║     ██╔══██║   ██║   ██╔══╝  
    ╚██████╔╝██║██║ ╚████║      ██║   ███████╗██║ ╚═╝ ██║██║     ███████╗██║  ██║   ██║   ███████╗
     ╚═════╝ ╚═╝╚═╝  ╚═══╝      ╚═╝   ╚══════╝╚═╝     ╚═╝╚═╝     ╚══════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝	

    欢迎使用 gin-template
	当前版本:V0.0.1 Beta
	默认后端接口运行地址:%s:%d
`, configs.Get().Project.Domain, configs.Get().Project.Port)

// 封装handler，统一处理
func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handler := handler
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

	mux := &mux{
		engine: gin.New(),
	}

	fmt.Println(color.Blue(_UI))

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
			pprof.Register(mux.engine) // register pprof to gin
		}
	}

	if !opt.disableSwagger {
		// swagger
		if !env.Active().IsPro() {
			mux.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // register swagger
		}
	}

	if !opt.disablePrometheus {
		// prometheus
		mux.engine.GET("/debug/metrics", gin.WrapH(promhttp.Handler())) // register prometheus
	}

	if opt.enableCors {
		// 跨域
		mux.engine.Use(MiddlewareCors())
	}

	if opt.enableOpenBrowser != "" {
		_ = browser.Open(opt.enableOpenBrowser)
	}

	// recover两次，防止处理时发生panic，尤其是在OnPanicNotify中。
	mux.engine.Use(MiddlewareRecover(logger))

	mux.engine.Use(MiddlewareTrace(logger, opt))

	// 限流
	if opt.enableRate {
		mux.engine.Use(MiddlewareLimit())
	}

	//
	mux.engine.NoMethod(wrapHandlers(DisableTraceLog)...)
	mux.engine.NoRoute(wrapHandlers(DisableTraceLog)...)

	system := mux.Group("/system")
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
				Status:      "ok777",
			}
			ctx.Payload(resp)
		})
	}

	return mux, nil
}
