package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"user-go/initializer"
	"user-go/web/handler"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func Init(service initializer.Service) {
	userHandler := handler.UserHandler{
		UserService: service.UserService,
	}
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.POST("/users", userHandler.Create)
	e.Logger.Fatal(e.Start(":8080"))
}
