package jwtHandlerClient

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/gommon/log"
	"user-go/domain/interfaces"
	"user-go/infra/jwthandler"
	"user-go/lib/authorization"
	"user-go/lib/unixtime"
)

type JwtHandlerClient struct {
	jwtHandler jwthandler.JwtHandler
	timekeeper interfaces.ITimeKeeper
}

var _ interfaces.IJwtHandlerClient = JwtHandlerClient{}

func NewJwtHandlerClient(
	jwtHandler jwthandler.JwtHandler,
	timekeeper interfaces.ITimeKeeper,
) JwtHandlerClient {
	return JwtHandlerClient{
		jwtHandler: jwtHandler,
		timekeeper: timekeeper,
	}
}

func (c JwtHandlerClient) Parse(rawToken authorization.TokenString) (authorization.Authorization, error) {
	parsedToken, err := c.jwtHandler.Parse(string(rawToken))
	if err != nil {
		return authorization.Authorization{}, err
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !(ok && parsedToken.Valid) {
		return authorization.Authorization{}, authorization.InvalidTokenString()
	}
	log.Print(claims)
	log.Print(claims)

	expInt64, ok := claims["exp"].(float64)
	if !(ok && unixtime.UnixTime(expInt64) >= c.timekeeper.Now()) {
		return authorization.Authorization{}, authorization.TokenStringExpired()
	}

	userID, ok := claims["userid"].(float64)
	if !ok {
		return authorization.Authorization{}, authorization.InvalidTokenString()
	}
	return authorization.NewAuthorization(int64(userID), authorization.TokenString(parsedToken.Raw)), nil
}
