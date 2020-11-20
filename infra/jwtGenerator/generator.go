package jwtGenerator

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtGenerator struct {
	key []byte
}

func NewJwtGenerator(key []byte) JwtGenerator {
	return JwtGenerator{
		key: key,
	}
}

func (g JwtGenerator) GenerateJwtToken(payload map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(payload))
	return token.SignedString(g.key)
}
