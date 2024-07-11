package tool

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type clearCacheRequest struct{}

type clearCacheResponse struct{}

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

	}
}
