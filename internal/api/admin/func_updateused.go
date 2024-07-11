package admin

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type updateUsedRequest struct{}

type updateUsedResponse struct{}

// UpdateUsed 更新管理员为启用/禁用
//
//	@Summary		更新管理员为启用/禁用
//	@Description	更新管理员为启用/禁用
//	@Tags			API.admin
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		updateUsedRequest	true	"请求信息"
//	@Success		200		{object}	updateUsedResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/admin/used [patch]
func (h *handler) UpdateUsed() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
