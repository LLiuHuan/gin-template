// Package core
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 17:42
package core

import (
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func MiddlewareLimit() gin.HandlerFunc {
	// TODO: 修改改成修改配置文件的方式
	limiter := rate.NewLimiter(rate.Every(time.Second*1), configs.MaxRequestsPerSecond)
	return func(ctx *gin.Context) {
		context := newContext(ctx)
		defer releaseContext(context)

		if !limiter.Allow() {
			context.AbortWithError(Error(
				http.StatusTooManyRequests,
				code.TooManyRequests,
				code.Text(code.TooManyRequests)),
			)
			return
		}

		ctx.Next()
	}
}
