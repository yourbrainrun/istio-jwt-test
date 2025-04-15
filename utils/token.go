package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
	"math/big"
	"time"
)

// load key
func loadPrivateKey() (*rsa.PrivateKey, error) {
	privKeyBytes, err := ioutil.ReadFile("./resource/cert/private-key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}
	return privateKeyInterface.(*rsa.PrivateKey), nil

}

// load public key
func loadPublicKey() (*rsa.PublicKey, error) {
	pubKeyBytes, err := ioutil.ReadFile("./resource/cert/public-key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse public key: %v", err)
	}
	return pubKeyInterface.(*rsa.PublicKey), nil

}

// GenerateJWT  JWT Token
func GenerateJWT() (string, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"sub":  "user123",
		"exp":  time.Now().Add(time.Hour * 48).Unix(),
		"iat":  time.Now().Unix(),
		"role": "admin",
		"iss":  "https://issuer1.com",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "rsa-key-id" // add kid column

	return token.SignedString(privateKey)
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	publicKey, err := loadPublicKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Base64 URL encode
func base64UrlEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// GenerateJWKS JWKS (JSON Web Key Set) JSON
func GenerateJWKS() map[string]interface{} {
	publicKey, err := loadPublicKey()
	if err != nil {
		return nil
	}

	kid := "rsa-key-id" // Key ID same to jwt

	n := base64UrlEncode(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())

	// generate JWKS JSON
	jwks := map[string]interface{}{
		"keys": []interface{}{
			map[string]interface{}{
				"kty": "RSA",
				"kid": kid,
				"alg": "RS256",
				"n":   n,
				"e":   e,
			},
		},
	}
	return jwks
}
