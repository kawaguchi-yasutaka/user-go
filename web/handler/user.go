package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"user-go/domain/model"
	"user-go/domain/service"
	"user-go/web/request"
)

type UserHandler struct {
	UserService service.UserService
}

func (handler UserHandler) Create(e echo.Context) error {
	req := request.UserCreateRequest{}
	if err := e.Bind(&req); err != nil {
		return err
	}

	email, err := model.NewUserEmail(req.Email)
	if err != nil {
		return err
	}
	password, err := model.NewUserRawPassword(req.Password)
	if err != nil {
		return err
	}

	if err := handler.UserService.Create(email, password); err != nil {
		return err
	}
	return e.NoContent(http.StatusNoContent)
}
