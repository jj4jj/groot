package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
)

//Encrypt encrypts the plain text with bcrypt
func EncryptPassword(source string) (string, error) {
	hashedBytes, e := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), e
}

// Compare compares the encrypted text with the plain text if it's the same
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoadRSAPrivateKey(privateKeyFile string) *rsa.PrivateKey {
	keybuffer, e := ioutil.ReadFile(privateKeyFile)
	if e != nil {
		log.Errorf(e, "loadPrivateKey error : %v", e)
		return nil
	}
	log.Debugf("the pri key buffer length is : %d", len(keybuffer))
	block, _ := pem.Decode([]byte(keybuffer))
	if block == nil {
		log.Errorf(nil, "loadPrivateKey decode error ")
		return nil
	}
	privateKey, e := x509.ParsePKCS8PrivateKey(block.Bytes)
	if e != nil {
		log.Errorf(e, "loadPrivateKey parse error ")
		return nil
	}
	resultPrivateKey := privateKey.(*rsa.PrivateKey)
	return resultPrivateKey
}

func LoadRSAPublicKey(publicKeyFile string) *rsa.PublicKey {
	keybuffer, e := ioutil.ReadFile(viper.GetString("key.rsa_public"))
	if e != nil {
		log.Errorf(e, "loadPublicKey error : %v", e)
		return nil
	}
	block, _ := pem.Decode([]byte(keybuffer))
	if block == nil {
		log.Errorf(nil, "loadPublicKey error : %v", e)
		return nil
	}
	publicKey, e := x509.ParsePKIXPublicKey(block.Bytes)
	if e != nil {
		log.Errorf(e, "loadPublicKey error : %v", e)
		return nil
	}
	resultPublicKey := publicKey.(*rsa.PublicKey)
	return resultPublicKey
}

func DecryptPasswordFromRSA(password string, privateKey *rsa.PrivateKey) (string, error) {
	decodetext, e := base64.StdEncoding.DecodeString(password)
	if e != nil {
		return "", e
	}
	decryptedText, e := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodetext)
	if e != nil {
		return "", e
	}
	return string(decryptedText), nil
}

func EncryptWithAES(b []byte, secret []byte) []byte {
	return b
}

func DecryptWithAES(b []byte, secret []byte) []byte {
	return b
}
