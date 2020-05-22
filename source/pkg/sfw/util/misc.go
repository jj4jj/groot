package util

import (
	"strings"
)

type RandCharsetFlag int

const (
	RAND_CHARSET_NUMBER RandCharsetFlag = 1 << iota
	RAND_CHARSET_ALPHA_LOWER
	RAND_CHARSET_ALPHA_UPPER
	RAND_CHARSET_SYMBOL
)

var (
	charsetNumber     []byte
	charsetAlphaLower []byte
	charsetAlphaUpper []byte
	charsetSymbol     []byte
	charsetNormal     []byte
)

func init() {
	charsetNumber = []byte("0123456789")
	strLower := "abcdefghijklmnopqrstuvwxyz"
	charsetAlphaLower = []byte(strLower)
	charsetAlphaUpper = []byte(strings.ToUpper(strLower))
	charsetSymbol = []byte("~!@#$%^&*()+-=_{}[]:';<>?,./`")
	charsetNormal = append(charsetNormal, charsetNumber...)
	charsetNormal = append(charsetNormal, charsetAlphaLower...)
	charsetNormal = append(charsetNormal, charsetAlphaUpper...)
}

func GetRandCharset(flag int) []byte {
	if flag == 0 {
		return charsetNormal
	}

	var result []byte
	if (flag & int(RAND_CHARSET_NUMBER)) != 0 {
		result = append(result, charsetNumber...)
	}
	if (flag & int(RAND_CHARSET_ALPHA_LOWER)) != 0 {
		result = append(result, charsetAlphaLower...)
	}
	if (flag & int(RAND_CHARSET_ALPHA_UPPER)) != 0 {
		result = append(result, charsetAlphaUpper...)
	}
	if (flag & int(RAND_CHARSET_SYMBOL)) != 0 {
		result = append(result, charsetSymbol...)
	}

	return result
}

func RandomAesKey() []byte {
	//16 or 32 bytes
	return []byte(RandomStringWithCharset(16, charsetNormal))
}

//RandomStr 随机生成字符串
func RandomString(length int) string {
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, charsetNormal[UtilRandom.Intn(len(charsetNormal))])
	}
	return string(result)
}

func RandomStringWithCharset(length int, charset []byte) string {
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, charset[UtilRandom.Intn(len(charset))])
	}
	return string(result)
}
