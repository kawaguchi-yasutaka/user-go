package web

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"user-go/initializer"
	"user-go/lib/myerror"
	"user-go/web/handler"
)

func customErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	errorType := "unknown"
	message := err.Error()
	if err, ok := err.(myerror.CustomError); ok {
		code = err.StatusCode
		errorType = string(err.ErrorType)
		message = err.Message
	}
	c.Logger().Error(err)
	c.JSON(code, map[string]interface{}{"error_message": message, "error_type": errorType})
}

func Init(service initializer.Service) {
	userHandler := handler.UserHandler{
		UserService: service.UserService,
	}
	e := echo.New()
	e.HTTPErrorHandler = customErrorHandler
	e.Use(middleware.Logger())
	e.POST("/users", userHandler.Create)
	e.GET("/activate_users", userHandler.Activate)
	e.POST("/login", userHandler.Login)
	e.POST("/logind", userHandler.Logind)
	e.Logger.Fatal(e.Start(":8080"))
}
