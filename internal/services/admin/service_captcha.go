// Package admin
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-23 22:47
package admin

import (
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/mojocn/base64Captcha"
	"github.com/spf13/cast"
	"time"
)

type CaptchaData struct {
	CaptchaId     string `json:"captchaId"`
	PicPath       string `json:"picPath"`
	CaptchaLength int    `json:"captchaLength"`
	OpenCaptcha   bool   `json:"openCaptcha"`
}

func (s *service) Captcha(ctx core.Context) (captchaData CaptchaData, err error) {
	openCaptcha := configs.Get().Captcha.OpenCaptcha
	openCaptchaTImeOut := configs.Get().Captcha.OpenCaptchaTimeOut
	key := ctx.GetCtx().ClientIP()
	v, err := s.cache.Get(key)
	if err != nil {
		err = s.cache.Set(key, "1", time.Duration(openCaptchaTImeOut)*time.Second)
		if err != nil {
			return CaptchaData{}, err
		}
	}

	var oc bool
	if openCaptcha == 0 || openCaptcha < cast.ToInt(v) {
		oc = true
	}

	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(configs.Get().Captcha.ImgHeight, configs.Get().Captcha.ImgWidth, configs.Get().Captcha.KeyLong, 0.7, 80)
	//cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(ctx.GetCtx())) // v8下使用redis
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		return CaptchaData{}, err
	}

	return CaptchaData{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: configs.Get().Captcha.KeyLong,
		OpenCaptcha:   oc,
	}, nil
}
