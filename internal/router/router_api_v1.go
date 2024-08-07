// Package router
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 15:32
package router

import (
	"github.com/LLiuHuan/gin-template/internal/api/admin"
	"github.com/LLiuHuan/gin-template/internal/api/helper"
	"github.com/LLiuHuan/gin-template/internal/api/tool"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

func setApiV1Router(r *resource) {
	// helper
	helperHandler := helper.New(r.logger, r.db, r.cache)
	toolHandler := tool.New(r.logger, r.db, r.cache)
	adminHandler := admin.New(r.logger, r.db, r.cache, r.dlp)

	apiRouter := r.mux.Group("/api/v1")
	{
		helpers := apiRouter.Group("/helper")
		{
			helpers.GET("/md5/:str", helperHandler.Md5())
			helpers.POST("/sign", helperHandler.Sign())
			helpers.POST("/upload", helperHandler.UploadFile())
			helpers.POST("/merge", helperHandler.UploadMerge())
			helpers.POST("/verify", helperHandler.UploadVerify())
		}

		notRBAC := apiRouter.Group("", core.DisableTraceLog, core.DisableRecordMetrics)
		{
			notRBAC.GET("/install")

			toolRouter := notRBAC.Group("/tool")
			{
				toolRouter.GET("/project/info", toolHandler.ProjectInfo())
				//	toolHandler := tool.New(r.logger, r.db, r.cache)
				toolRouter.GET("/hashids/encode/:id", core.AliasForRecordMetrics("/api/tool/hashids/encode"), toolHandler.HashIdsEncode())
				toolRouter.GET("/hashids/decode/:id", core.AliasForRecordMetrics("/api/tool/hashids/decode"), toolHandler.HashIdsDecode())
				toolRouter.POST("/cache/search", toolHandler.SearchCache())
				toolRouter.PATCH("/cache/clear", toolHandler.ClearCache())
				toolRouter.GET("/project/logs", toolHandler.Logs())
				//	api.GET("/tool/data/dbs", toolHandler.Dbs())
				//	api.POST("/tool/data/tables", toolHandler.Tables())
				//	api.POST("/tool/data/mysql", toolHandler.SearchMySQL())
				//	api.POST("/tool/send_message", toolHandler.SendMessage())
			}
		}

		//api := apiRouter.Group("", core.WrapAuthHandler(r.interceptors.CheckLogin), r.interceptors.CheckSignature(), r.interceptors.CheckRBAC())
		//{
		//	toolRouter := api.Group("/tool")
		//	{
		//		toolRouter.GET("/project/info", toolHandler.ProjectInfo())
		//	}
		//}
	}
	// 需要签名验证，无需登录验证，无需 RBAC 权限验证
	login := r.mux.Group("/api/v1", r.interceptors.CheckSignature())
	{
		login.POST("/login", adminHandler.Login())
		login.POST("/captcha", adminHandler.Captcha())
	}

	// 需要签名验证、登录验证，无需 RBAC 权限验证
	notRBAC := r.mux.Group("/api/v1", core.WrapAuthHandler(r.interceptors.CheckLogin), r.interceptors.CheckSignature())
	{
		//notRBAC.POST("/admin/logout", adminHandler.Logout())
		//notRBAC.PATCH("/admin/modify_password", adminHandler.ModifyPassword())
		notRBAC.GET("/admin/info", adminHandler.Detail())
		//notRBAC.PATCH("/admin/modify_personal_info", adminHandler.ModifyPersonalInfo())
	}

	//// 需要签名验证、登录验证、RBAC 权限验证
	//api := r.mux.Group("/api", core.WrapAuthHandler(r.interceptors.CheckLogin), r.interceptors.CheckSignature(), r.interceptors.CheckRBAC())
	//{
	//	// authorized
	//	authorizedHandler := authorized.New(r.logger, r.db, r.cache)
	//	api.POST("/authorized", authorizedHandler.Create())
	//	api.GET("/authorized", authorizedHandler.List())
	//	api.PATCH("/authorized/used", authorizedHandler.UpdateUsed())
	//	api.DELETE("/authorized/:id", core.AliasForRecordMetrics("/api/authorized/info"), authorizedHandler.Delete())
	//
	//	api.POST("/authorized_api", authorizedHandler.CreateAPI())
	//	api.GET("/authorized_api", authorizedHandler.ListAPI())
	//	api.DELETE("/authorized_api/:id", core.AliasForRecordMetrics("/api/authorized_api/info"), authorizedHandler.DeleteAPI())
	//
	//	api.POST("/admin", adminHandler.Create())
	//	api.GET("/admin", adminHandler.List())
	//	api.PATCH("/admin/used", adminHandler.UpdateUsed())
	//	api.PATCH("/admin/offline", adminHandler.Offline())
	//	api.PATCH("/admin/reset_password/:id", core.AliasForRecordMetrics("/api/admin/reset_password"), adminHandler.ResetPassword())
	//	api.DELETE("/admin/:id", core.AliasForRecordMetrics("/api/admin"), adminHandler.Delete())
	//
	//	api.POST("/admin/menu", adminHandler.CreateAdminMenu())
	//	api.GET("/admin/menu/:id", core.AliasForRecordMetrics("/api/admin/menu"), adminHandler.ListAdminMenu())
	//
	//	// menu
	//	menuHandler := menu.New(r.logger, r.db, r.cache)
	//	api.POST("/menu", menuHandler.Create())
	//	api.GET("/menu", menuHandler.List())
	//	api.GET("/menu/:id", core.AliasForRecordMetrics("/api/menu"), menuHandler.Detail())
	//	api.PATCH("/menu/used", menuHandler.UpdateUsed())
	//	api.PATCH("/menu/sort", menuHandler.UpdateSort())
	//	api.DELETE("/menu/:id", core.AliasForRecordMetrics("/api/menu"), menuHandler.Delete())
	//	api.POST("/menu_action", menuHandler.CreateAction())
	//	api.GET("/menu_action", menuHandler.ListAction())
	//	api.DELETE("/menu_action/:id", core.AliasForRecordMetrics("/api/menu_action"), menuHandler.DeleteAction())
	//
	//	// tool
	//	toolHandler := tool.New(r.logger, r.db, r.cache)
	//	api.GET("/tool/hashids/encode/:id", core.AliasForRecordMetrics("/api/tool/hashids/encode"), toolHandler.HashIdsEncode())
	//	api.GET("/tool/hashids/decode/:id", core.AliasForRecordMetrics("/api/tool/hashids/decode"), toolHandler.HashIdsDecode())
	//	api.POST("/tool/cache/search", toolHandler.SearchCache())
	//	api.PATCH("/tool/cache/clear", toolHandler.ClearCache())
	//	api.GET("/tool/data/dbs", toolHandler.Dbs())
	//	api.POST("/tool/data/tables", toolHandler.Tables())
	//	api.POST("/tool/data/mysql", toolHandler.SearchMySQL())
	//	api.POST("/tool/send_message", toolHandler.SendMessage())
	//
	//	// config
	//	configHandler := config.New(r.logger, r.db, r.cache)
	//	api.PATCH("/config/email", configHandler.Email())
	//
	//	// cron
	//	cronHandler := cron.New(r.logger, r.db, r.cache, r.cronServer)
	//	api.POST("/cron", cronHandler.Create())
	//	api.GET("/cron", cronHandler.List())
	//	api.GET("/cron/:id", core.AliasForRecordMetrics("/api/cron/detail"), cronHandler.Detail())
	//	api.POST("/cron/:id", core.AliasForRecordMetrics("/api/cron/modify"), cronHandler.Modify())
	//	api.PATCH("/cron/used", cronHandler.UpdateUsed())
	//	api.PATCH("/cron/exec/:id", core.AliasForRecordMetrics("/api/cron/exec"), cronHandler.Execute())
	//
	//}
}
