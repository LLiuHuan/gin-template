// Package interceptor
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 15:00
// @description: 登录拦截器
package interceptor

import (
	"encoding/json"
	"net/http"

	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/proposal"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/LLiuHuan/gin-template/pkg/errors"
)

func (i *interceptor) CheckLogin(ctx core.Context) (sessionUserInfo proposal.SessionUserInfo, err core.BusinessError) {
	token := ctx.GetHeader(configs.HeaderLoginToken)
	if token == "" {
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(errors.New("Header 中缺少 Token 参数"))

		return
	}

	if !i.cache.Exists(configs.RedisKeyPrefixLoginUser + token) {
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(errors.New("请先登录"))

		return
	}

	cacheData, cacheErr := i.cache.Get(configs.RedisKeyPrefixLoginUser+token, redis.WithTrace(ctx.Trace()))
	if cacheErr != nil {
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(cacheErr)

		return
	}

	jsonErr := json.Unmarshal([]byte(cacheData), &sessionUserInfo)
	if jsonErr != nil {
		core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(jsonErr)

		return
	}

	return
}
