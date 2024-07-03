// Package core
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 14:40
package core

import "github.com/LLiuHuan/gin-template/internal/proposal"

type Option func(*option)

type option struct {
	disablePProf      bool
	disableSwagger    bool
	disablePrometheus bool
	enableCors        bool
	enableRate        bool
	enableOpenBrowser string
	alertNotify       proposal.NotifyHandler
	recordHandler     proposal.RecordHandler
}

// WithDisablePProf 禁用 pprof
func WithDisablePProf() Option {
	return func(opt *option) {
		opt.disablePProf = true
	}
}

// WithDisableSwagger 禁用 swagger
func WithDisableSwagger() Option {
	return func(opt *option) {
		opt.disableSwagger = true
	}
}

// WithDisablePrometheus 禁用prometheus
func WithDisablePrometheus() Option {
	return func(opt *option) {
		opt.disablePrometheus = true
	}
}

// WithAlertNotify 设置告警通知
func WithAlertNotify(notifyHandler proposal.NotifyHandler) Option {
	return func(opt *option) {
		opt.alertNotify = notifyHandler
	}
}

// WithRecordMetrics 设置记录接口指标
func WithRecordMetrics(recordHandler proposal.RecordHandler) Option {
	return func(opt *option) {
		opt.recordHandler = recordHandler
	}
}

// WithEnableOpenBrowser 启动后在浏览器中打开 uri
func WithEnableOpenBrowser(uri string) Option {
	return func(opt *option) {
		opt.enableOpenBrowser = uri
	}
}

// WithEnableCors 设置支持跨域
func WithEnableCors() Option {
	return func(opt *option) {
		opt.enableCors = true
	}
}

// WithEnableRate 设置支持限流
func WithEnableRate() Option {
	return func(opt *option) {
		opt.enableRate = true
	}
}

// DisableTraceLog 禁止记录日志
func DisableTraceLog(ctx Context) {
	ctx.disableTrace()
}

// DisableRecordMetrics 禁止记录指标
func DisableRecordMetrics(ctx Context) {
	ctx.disableRecordMetrics()
}

// AliasForRecordMetrics 对请求路径起个别名，用于记录指标。
// 如：Get /user/:username 这样的路径，因为 username 会有非常多的情况，这样记录指标非常不友好。
func AliasForRecordMetrics(path string) HandlerFunc {
	return func(ctx Context) {
		ctx.setAlias(path)
	}
}

// WrapAuthHandler 用来处理 Auth 的入口
func WrapAuthHandler(handler func(Context) (sessionUserInfo proposal.SessionUserInfo, err BusinessError)) HandlerFunc {
	return func(ctx Context) {
		sessionUserInfo, err := handler(ctx)
		if err != nil {
			ctx.AbortWithError(err)
			return
		}

		ctx.setSessionUserInfo(sessionUserInfo)
	}
}
