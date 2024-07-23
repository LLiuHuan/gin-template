// Package core
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-19 16:14
package core

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"go.uber.org/zap"
)

func MiddlewareTls(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:443",
		})
		err := middleware.Process(ctx.Writer, ctx.Request)
		if err != nil {
			// 如果出现错误，请不要继续
			logger.Error("tls middleware error", zap.Error(err))
			//ctx.Abort()
			return
		}
		// 如果响应是重定向，请避免标头重写。
		if status := ctx.Writer.Status(); status > 300 && status < 399 {
			ctx.Abort()
		}
		// 继续往下处理
		ctx.Next()
	}
}
