package tool

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type hashIdsDecodeRequest struct{}

type hashIdsDecodeResponse struct{}

// HashIdsDecode HashIds 解密
//
//	@Summary		HashIds 解密
//	@Description	HashIds 解密
//	@Tags			API.tool
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		hashIdsDecodeRequest	true	"请求信息"
//	@Success		200		{object}	hashIdsDecodeResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/v1/tool/hashids/decode/{id} [get]
func (h *handler) HashIdsDecode() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
