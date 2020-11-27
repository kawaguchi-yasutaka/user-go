package jwthandler

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type JwtHandler struct {
	key *rsa.PublicKey
}

func NewJwtHandler(key []byte) JwtHandler {
	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		panic(err)
	}
	return JwtHandler{
		key: parsedKey,
	}
}

func (h JwtHandler) Parse(rawToken string) (*jwt.Token, error) {
	return jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return h.key, nil
	})
}
