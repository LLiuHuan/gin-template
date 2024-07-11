package admin

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type listAdminMenuRequest struct{}

type listAdminMenuResponse struct{}

// ListAdminMenu 菜单授权列表
//
//	@Summary		菜单授权列表
//	@Description	菜单授权列表
//	@Tags			API.admin
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		listAdminMenuRequest	true	"请求信息"
//	@Success		200		{object}	listAdminMenuResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/admin/menu/{id} [get]
func (h *handler) ListAdminMenu() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
