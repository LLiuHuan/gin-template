// Package router
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-08-22 03:02
// @description: render 路由
package router

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/render/cron"
	"github.com/LLiuHuan/gin-template/internal/render/generator"
	"github.com/LLiuHuan/gin-template/internal/render/install"
)

func setRenderRouter(r *resource) {
	renderInstall := install.New(r.logger)
	renderGenerator := generator.New(r.logger, r.db, r.cache)
	renderCron := cron.New(r.logger, r.db, r.cache)

	// 无需记录日志，无需 RBAC 权限验证
	notRBAC := r.mux.Group("", core.DisableTraceLog, core.DisableRecordMetrics)
	{
		// 安装
		notRBAC.GET("/install", renderInstall.View())
		notRBAC.POST("/install/execute", renderInstall.Execute())
	}

	// 无需记录日志，需要 RBAC 权限验证
	render := r.mux.Group("", core.DisableTraceLog, core.DisableRecordMetrics)
	{
		// 代码生成器
		render.GET("/generator/gorm", renderGenerator.GormView())
		render.POST("/generator/gorm/execute", renderGenerator.GormExecute())

		render.GET("/generator/handler", renderGenerator.HandlerView())
		render.POST("/generator/handler/execute", renderGenerator.HandlerExecute())

		// 后台任务
		render.GET("/cron/list", renderCron.List())
		render.GET("/cron/add", renderCron.Add())
		render.GET("/cron/edit/:id", renderCron.Edit())
	}
}
