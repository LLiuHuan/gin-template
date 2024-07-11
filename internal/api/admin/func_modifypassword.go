package admin

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type modifyPasswordRequest struct{}

type modifyPasswordResponse struct{}

// ModifyPassword 修改密码
//
//	@Summary		修改密码
//	@Description	修改密码
//	@Tags			API.admin
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		modifyPasswordRequest	true	"请求信息"
//	@Success		200		{object}	modifyPasswordResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/admin/modify_password [patch]
func (h *handler) ModifyPassword() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
