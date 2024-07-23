package tool

import (
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"net/http"
)

type hashIdsEncodeRequest struct {
	Id uint64 `uri:"id"` // 需加密的数字
}

type hashIdsEncodeResponse struct {
	Val string `json:"val"` // 加密后的值
}

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
		req := new(hashIdsEncodeRequest)
		res := new(hashIdsEncodeResponse)
		if err := ctx.ShouldBindURI(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		hashId, err := h.hashids.HashidsEncode([]uint64{req.Id})
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.HashIdsEncodeError,
				code.Text(code.HashIdsEncodeError)).WithError(err),
			)
			return
		}

		res.Val = hashId

		ctx.Payload(res)
	}
}
