package admin

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type logoutRequest struct{}

type logoutResponse struct{}

// Logout 管理员登出
//
//	@Summary		管理员登出
//	@Description	管理员登出
//	@Tags			API.admin
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		logoutRequest	true	"请求信息"
//	@Success		200		{object}	logoutResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/admin/logout [post]
func (h *handler) Logout() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
