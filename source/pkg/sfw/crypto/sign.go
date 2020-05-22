package crypto

import "crypto/rsa"

var (
	signPrivateKey *rsa.PrivateKey
)

func SetPrivateKey(k *rsa.PrivateKey) {
	signPrivateKey = k
}

func GenerateSha256WithRsa(b []byte) string {
	//todo
	return ""
}
func VerifySha256WithRsa(b []byte, sign string) bool {
	//todo
	return false
}
