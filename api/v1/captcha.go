package v1

import (
	"gin-blog/model"

	"github.com/gin-gonic/gin"
)

type BaseApi struct{}

// 生成验证码
func (*BaseApi) Captcha(c *gin.Context) {
	// 判断验证是否开启
	captchaResp := model.CaptchaResponse{
		CaptchaId:     "123",
		PicPath:       "b64s",
		CaptchaLength: 6,
		OpenCaptcha:   true,
	}
	c.JSON(200, captchaResp)
}

// {
// 	CaptchaId:     id,
// 	PicPath:       b64s,
// 	CaptchaLength: global.GVA_CONFIG.Captcha.KeyLong,
// 	OpenCaptcha:   oc,
// }
