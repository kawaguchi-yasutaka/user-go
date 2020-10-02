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

func (handler UserHandler) Create(c echo.Context) error {
	req := request.UserCreateRequest{}
	if err := c.Bind(&req); err != nil {
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

	userId, err := handler.UserService.Create(email, password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"id": userId})
}

func (handler UserHandler) Activate(c echo.Context) error {
	if err := handler.UserService.Activate(
		model.UserActivationCode(c.QueryParam("code")),
	); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

//dbにあるか
//期限内か
//すでに認証済である。
//activateに更新する。
