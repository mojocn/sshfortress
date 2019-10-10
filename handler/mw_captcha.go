package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

//MwCaptchaCheck 图形验证码中间件
func MwCaptchaCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := c.Query("captcha_id")
		cVal := c.Query("captcha_value")
		if cid == "" || cVal == "" {
			jsonError(c, fmt.Sprintf("URL参数 %s %s 不能为空", "captcha_id", "captcha_value"))
			return
		}
		//万能图形验证码
		if base64Captcha.VerifyCaptcha(cid, cVal) {
			c.Next()
		} else {
			jsonError(c, "验证码错误")
		}
	}
}
