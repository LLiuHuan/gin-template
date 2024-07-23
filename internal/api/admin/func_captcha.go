// Package admin
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-23 22:32
package admin

import (
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"net/http"
)

type captchaRequest struct{}

type captchaResponse struct {
	CaptchaId     string `json:"captchaId"`
	PicPath       string `json:"picPath"`
	CaptchaLength int    `json:"captchaLength"`
	OpenCaptcha   bool   `json:"openCaptcha"`
}

// Captcha 获取验证码
//
//	@Summary		获取验证码
//	@Description	获取验证码
//	@Tags			API.admin
//	@Accept			application/json
//	@Produce		json
//	@Param			Request	body		captchaRequest	true	"请求信息"
//	@Success		200		{object}	captchaResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/admin/captcha [post]
func (h *handler) Captcha() core.HandlerFunc {
	return func(ctx core.Context) {
		res := new(captchaResponse)
		data, err := h.adminService.Captcha(ctx)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminCaptchaError,
				code.Text(code.AdminCaptchaError)).WithError(err),
			)
			return
		}

		res.CaptchaId = data.CaptchaId
		res.PicPath = data.PicPath
		res.CaptchaLength = data.CaptchaLength
		res.OpenCaptcha = data.OpenCaptcha

		ctx.Payload(res)
	}
}
