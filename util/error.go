package util

type CaptchaNotFound struct{}

func (m *CaptchaNotFound) Error() string {
	return "Captcha not found"
}
