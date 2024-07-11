package tool

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type searchCacheRequest struct{}

type searchCacheResponse struct{}

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

	}
}
