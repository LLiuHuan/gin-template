package admin

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type detailRequest struct{}

type detailResponse struct{}

// Detail 个人信息
//
//	@Summary		个人信息
//	@Description	个人信息
//	@Tags			API.admin
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		detailRequest	true	"请求信息"
//	@Success		200		{object}	detailResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/admin/info [get]
func (h *handler) Detail() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
