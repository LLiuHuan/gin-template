package admin

import (
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/services/admin"
	"net/http"
)

type detailRequest struct{}

type detailResponse struct {
	Username string                 `json:"username" mask:"NAME"`     // 用户名
	Nickname string                 `json:"nickname" mask:"NAME"`     // 昵称
	Mobile   string                 `json:"mobile" mask:"CHINAPHONE"` // 手机号
	Menu     []admin.ListMyMenuData `json:"menu"`                     // 菜单栏
}

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
		res := new(detailResponse)

		searchOneData := new(admin.SearchOneData)
		searchOneData.Id = ctx.SessionUserInfo().UserID
		searchOneData.IsUsed = 1

		info, err := h.adminService.Detail(ctx, searchOneData)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminDetailError,
				code.Text(code.AdminDetailError)).WithError(err),
			)
			return
		}

		res.Username = info.Username
		res.Nickname = info.Nickname
		res.Mobile = info.Mobile
		outRes, err := h.dlp.MaskStruct(res)
		if err != nil {
			ctx.Payload(res)
			return
		}
		ctx.Payload(outRes)
	}
}
