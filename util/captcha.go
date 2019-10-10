package util

import (
	"github.com/mojocn/base64Captcha"
	"math/rand"
)

//demo http://captcha.mojotv.cn/
//github.com/mojocn/base64Captcha
var captchaOption = base64Captcha.ConfigCharacter{
	Height:             39,
	Width:              210,
	CaptchaLen:         4,
	ComplexOfNoiseDot:  1,
	ComplexOfNoiseText: 2,
	IsShowHollowLine:   false,
	IsShowNoiseDot:     false,
	IsShowNoiseText:    false,
	IsShowSineLine:     false,
	IsShowSlimeLine:    false,
	IsUseSimpleFont:    true,
	Mode:               2,
}

//GetCaptchaImage 产生base64的图形验证码
//使用默认的单机存储

//如果需要多台机器请自定义store
//如果lvs 部署定义redis_store  https://github.com/mojocn/base64Captcha/blob/master/_examples_redis/main.go
func GetCaptchaImage() (id, b64image string) {
	captchaOption.Mode = rand.Intn(4)
	captchaId, captchaInterfaceInstance := base64Captcha.GenerateCaptcha("", captchaOption)
	base64blob := base64Captcha.CaptchaWriteToBase64Encoding(captchaInterfaceInstance)
	return captchaId, base64blob
}
