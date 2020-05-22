package constk

type CaptchaSceneType uint

const (
	CAPTCHA_ACCOUNT_REGISTER CaptchaSceneType = 1001 + iota
	CAPTCHA_ACCOUNT_LOGIN
)
