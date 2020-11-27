package jwtGenerator

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
)

type JwtGenerator struct {
	key *rsa.PrivateKey
}

func NewJwtGenerator(key []byte) JwtGenerator {
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		panic(err)
	}
	return JwtGenerator{
		key: parsedKey,
	}
}

func (g JwtGenerator) Generate(claim map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claim))
	return token.SignedString(g.key)
}
