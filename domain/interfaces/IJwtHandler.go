package interfaces

import "github.com/dgrijalva/jwt-go"

type IJwtHandler interface {
	Parse(rawToken string) (*jwt.Token, error)
}
