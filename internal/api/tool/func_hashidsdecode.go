package tool

import (
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"net/http"
)

type hashIdsDecodeRequest struct {
	Id string `uri:"id"` // 需解密的密文
}

type hashIdsDecodeResponse struct {
	Val uint64 `json:"val"` // 解密后的值
}

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
		req := new(hashIdsDecodeRequest)
		res := new(hashIdsDecodeResponse)
		if err := ctx.ShouldBindURI(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		hashId, err := h.hashids.HashidsDecode(req.Id)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.HashIdsDecodeError,
				code.Text(code.HashIdsDecodeError)).WithError(err),
			)
			return
		}

		res.Val = hashId[0]

		ctx.Payload(res)
	}
}
