package middlewares

import "github.com/labstack/echo/v4"

var AuthInfo = "AuthInfo"

func AuthorizationMiddleware(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		c.Set(AuthInfo, auth)
		return f(c)
	}
}

func GetToken(c echo.Context) string {
	authInfo := c.Get(AuthInfo)
	if s, ok := authInfo.(string); ok {
		return s
	}
	return ""
}
