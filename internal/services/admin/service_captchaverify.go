// Package admin
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-23 23:38
package admin

import (
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/spf13/cast"
	"time"
)

func (s *service) CaptchaVerify(ctx core.Context, captcha string, captchaId string) bool {
	openCaptcha := configs.Get().Captcha.OpenCaptcha
	openCaptchaTImeOut := configs.Get().Captcha.OpenCaptchaTimeOut
	key := ctx.GetCtx().ClientIP()

	v, err := s.cache.Get(key)
	if err != nil {
		err = s.cache.Set(key, "1", time.Duration(openCaptchaTImeOut)*time.Second)
		if err != nil {
			return false
		}
	}

	var oc = openCaptcha == 0 || openCaptcha < cast.ToInt(v)

	if !oc || (captcha != "" && captchaId != "" && store.Verify(captchaId, captcha, true)) {
		return true
	}

	s.cache.Incr(key)

	return false
}
