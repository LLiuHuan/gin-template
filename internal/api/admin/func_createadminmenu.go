package admin

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type createAdminMenuRequest struct{}

type createAdminMenuResponse struct{}

// CreateAdminMenu 提交菜单授权
//
//	@Summary		提交菜单授权
//	@Description	提交菜单授权
//	@Tags			API.admin
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		createAdminMenuRequest	true	"请求信息"
//	@Success		200		{object}	createAdminMenuResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/admin/menu [post]
func (h *handler) CreateAdminMenu() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
