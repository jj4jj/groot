package account

import (
	"crypto/rsa"
	"errors"
	"groot/sfw/crypto"
)

var (
	authPrivateKey *rsa.PrivateKey
	authPublicKey  *rsa.PublicKey
)

func LoadAuthRSAPrivateKey(pkpath string) error {
	if authPrivateKey = crypto.LoadRSAPrivateKey(pkpath); authPrivateKey != nil {
		return errors.New("load private key error !")
	}
	return nil
}
func LoadAuthRSAPublicKey(pkpath string) error {
	if authPublicKey = crypto.LoadRSAPublicKey(pkpath); authPublicKey != nil {
		return errors.New("load pubkey error !")
	}
	return nil
}

func SafeCheckAuthPass(reqPass, dbPass string) bool {
	decPass, e := crypto.DecryptPasswordFromRSA(reqPass, authPrivateKey)
	if e != nil {
		return false
	}
	if e := crypto.ComparePassword(dbPass, decPass); e != nil {
		return false
	}
	return true
}
