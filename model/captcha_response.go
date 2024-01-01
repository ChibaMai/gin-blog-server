package model

type CaptchaResponse struct {
	CaptchaId     string `json:"CaptchaId"`
	PicPath       string `json:"PicPath"`
	CaptchaLength int    `json:"CaptchaLength"`
	OpenCaptcha   bool   `json:"OpenCaptcha"`
}
