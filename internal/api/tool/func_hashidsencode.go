package tool

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type hashIdsEncodeRequest struct{}

type hashIdsEncodeResponse struct{}

// HashIdsEncode HashIds 加密
//
//	@Summary		HashIds 加密
//	@Description	HashIds 加密
//	@Tags			API.tool
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		hashIdsEncodeRequest	true	"请求信息"
//	@Success		200		{object}	hashIdsEncodeResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/v1/tool/hashids/encode/{id} [get]
func (h *handler) HashIdsEncode() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
