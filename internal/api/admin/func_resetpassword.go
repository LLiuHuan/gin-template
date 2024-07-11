package admin

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type resetPasswordRequest struct{}

type resetPasswordResponse struct{}

// ResetPassword 重置密码
//
//	@Summary		重置密码
//	@Description	重置密码
//	@Tags			API.admin
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		resetPasswordRequest	true	"请求信息"
//	@Success		200		{object}	resetPasswordResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/admin/reset_password/{id} [patch]
func (h *handler) ResetPassword() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
