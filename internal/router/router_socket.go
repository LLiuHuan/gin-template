// Package router
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-04 17:33
package router

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/websocket/install"
)

func setSocketRouter(r *resource) {
	installSocket := install.New(r.logger, r.db, r.cache)

	// 无需记录日志
	// 无需记录日志
	socket := r.mux.Group("/socket", core.DisableTraceLog, core.DisableRecordMetrics)
	{
		// 系统消息
		socket.GET("/install", installSocket.Install())
	}
}
