package tool

import (
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"net/http"
)

type searchCacheRequest struct {
	RedisKey string `form:"redis_key"` // Redis Key
}

type searchCacheResponse struct {
	Val string `json:"val"` // 查询后的值
	TTL string `json:"ttl"` // 过期时间
}

// SearchCache 查询缓存
//
//	@Summary		查询缓存
//	@Description	查询缓存
//	@Tags			API.tool
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		searchCacheRequest	true	"请求信息"
//	@Success		200		{object}	searchCacheResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/v1/tool/cache/search [post]
func (h *handler) SearchCache() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(searchCacheRequest)
		res := new(searchCacheResponse)
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
		val, err := h.cache.Get(req.RedisKey, redis.WithTrace(ctx.Trace()))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.CacheGetError,
				code.Text(code.CacheGetError)).WithError(err),
			)
			return
		}

		ttl, err := h.cache.TTL(req.RedisKey)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.CacheGetError,
				code.Text(code.CacheGetError)).WithError(err),
			)
			return
		}

		res.Val = val
		res.TTL = ttl.String()
		ctx.Payload(res)
	}
}
