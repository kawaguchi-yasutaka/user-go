package jwtgenerator

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"user-go/domain/model"
	"user-go/lib/unixtime"
)

type JwtGenerator struct {
	key *rsa.PrivateKey
}

type MyCustomClaim struct {
	UserId int64 `json:"userid"`
	jwt.StandardClaims
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

func (g JwtGenerator) Generate(userId model.UserID, exp unixtime.UnixTime) (string, error) {
	claims := MyCustomClaim{
		int64(userId),
		jwt.StandardClaims{
			ExpiresAt: int64(exp),
			Issuer:    "user-go",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(g.key)
}
