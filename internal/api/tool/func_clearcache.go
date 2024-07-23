package tool

import (
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"net/http"
)

type clearCacheRequest struct {
	RedisKey string `form:"redis_key"` // Redis Key

}

type clearCacheResponse struct {
	Bool bool `json:"bool"` // 删除结果
}

// ClearCache 清空缓存
//
//	@Summary		清空缓存
//	@Description	清空缓存
//	@Tags			API.tool
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		clearCacheRequest	true	"请求信息"
//	@Success		200		{object}	clearCacheResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/v1/tool/cache/clear [patch]
func (h *handler) ClearCache() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(clearCacheRequest)
		res := new(clearCacheResponse)
		if err := ctx.ShouldBindForm(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		if b := h.cache.Exists(req.RedisKey); b != true {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.CacheNotExist,
				code.Text(code.CacheNotExist)),
			)
			return
		}

		b := h.cache.Del(req.RedisKey, redis.WithTrace(ctx.Trace()))
		if b != true {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.CacheDelError,
				code.Text(code.CacheDelError)),
			)
			return
		}

		res.Bool = b
		ctx.Payload(res)
	}
}
