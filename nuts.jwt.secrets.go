package gonuts

import (
	"crypto/rsa"
	"io/ioutil"
	"os"

	"github.com/golang-jwt/jwt"
)

var JWT_AccessToken_SignKey *rsa.PrivateKey
var JWT_AccessToken_VerifyKey *rsa.PublicKey
var JWT_RefreshToken_Secret string

// var JWT_AccessToken_SignKeyString []byte

func InitSecrets() {
	secretDir, _ := os.Getwd()
	JWT_AccessToken_SignKey = ReadFilePrivateKey(secretDir + "/secrets/jwt_xs_rs256.key")
	JWT_AccessToken_VerifyKey = ReadFilePublicKey(secretDir + "/secrets/jwt_xs_rs256.key.pub")
	// data, err := ioutil.ReadFile(secretDir + "/secrets/singleline.key")
	// data, err := ioutil.ReadFile(secretDir + "/secrets/jwt_xs_rs256.key")
	// if err != nil {
	// 	L.Errorf("[nuts.InitSecrets] %s", err)
	// 	os.Exit(1)
	// }
	// JWT_AccessToken_SignKeyString = data

}

func ReadFilePrivateKey(filepath string) *rsa.PrivateKey {
	signBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		L.Errorf("[nuts.ReadFilePrivateKey] FATAL ERROR-> %s", err)
		panic(err)
	}
	privateKeyImported, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		L.Errorf("[nuts.ReadFilePrivateKey] FATAL ERROR-> %s", err)
		panic(err)
	}
	JWT_RefreshToken_Secret = string(signBytes)[100:200]
	// L.Debugf("::::::::::::::::::::: JWT_RefreshToken_Secret = [%s]", JWT_RefreshToken_Secret)
	return privateKeyImported
}

func ReadFilePublicKey(filepath string) *rsa.PublicKey {
	signBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		L.Errorf("[nuts.ReadFilePublicKey] FATAL ERROR-> %s", err)
		panic(err)
	}
	privateKeyImported, err := jwt.ParseRSAPublicKeyFromPEM(signBytes)
	if err != nil {
		L.Errorf("[nuts.ReadFilePublicKey] FATAL ERROR-> %s", err)
		panic(err)
	}
	return privateKeyImported
}
