// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-23 22:36
package configs

type Captcha struct {
	KeyLong            int `toml:"keyLong"`            // KeyLong 验证码长度
	ImgWidth           int `toml:"imgWidth"`           // ImgWidth 图片宽度
	ImgHeight          int `toml:"imgHeight"`          // ImgHeight 图片高度
	OpenCaptcha        int `toml:"openCaptcha"`        // OpenCaptcha 防爆破验证码开启此数，0代表每次登录都需要验证码，其他数字代表错误密码此数，如3代表错误三次后出现验证码
	OpenCaptchaTimeOut int `toml:"openCaptchaTimeout"` // OpenCaptchaTimeOut 防爆破验证码超时时间，单位：s(秒)
}
