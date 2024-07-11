package admin

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type offlineRequest struct{}

type offlineResponse struct{}

// Offline 下线管理员
//
//	@Summary		下线管理员
//	@Description	下线管理员
//	@Tags			API.admin
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		offlineRequest	true	"请求信息"
//	@Success		200		{object}	offlineResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/admin/offline [patch]
func (h *handler) Offline() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
