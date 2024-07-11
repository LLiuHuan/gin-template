package admin

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type modifyPersonalInfoRequest struct{}

type modifyPersonalInfoResponse struct{}

// ModifyPersonalInfo 修改个人信息
//
//	@Summary		修改个人信息
//	@Description	修改个人信息
//	@Tags			API.admin
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		modifyPersonalInfoRequest	true	"请求信息"
//	@Success		200		{object}	modifyPersonalInfoResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/admin/modify_personal_info [patch]
func (h *handler) ModifyPersonalInfo() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
