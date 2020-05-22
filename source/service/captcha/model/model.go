package model

type (
	DbCaptchaState struct {
		Uuid         string `gorm:"UNIQUE_INDEX"` //for code key
		Code         string
		Target       string `gorm:"UNIQUE_INDEX:captcha_uniq_id"`
		Scene        uint32 `gorm:"UNIQUE_INDEX:captcha_uniq_id"`
		Data         string
		CheckedTimes int
		LastGenTime  int64
		ExpiredTime  int64
	}
)

