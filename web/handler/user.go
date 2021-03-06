package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"user-go/domain/model"
	"user-go/domain/service"
	"user-go/lib/authorization"
	"user-go/web/middlewares"
	"user-go/web/request"
	"user-go/web/response"
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	if err := handler.UserService.Activate(
		model.UserActivationCode(c.QueryParam("code")),
		model.UserID(id),
	); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (handler UserHandler) Login(c echo.Context) error {
	req := request.UserLoginRequest{}
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

	sessionId, err := handler.UserService.Login(email, password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sessionId": sessionId,
	})
}

func (handler UserHandler) Logind(c echo.Context) error {
	token := middlewares.GetToken(c)
	_, err := handler.UserService.Logind(authorization.TokenString(token))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (handler UserHandler) MultiAuthenticate(c echo.Context) error {
	req := request.UserMultiAuthenticateRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := handler.UserService.MultiAuthenticate(
		model.UserMultiAuthenticationCode(c.QueryParam("code")),
		model.UserSessionId(req.SessionId),
	); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (handler UserHandler) MultiAuthenticateAndGetJWT(c echo.Context) error {
	req := request.UserMultiAuthenticateRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	token, err := handler.UserService.MultiAuthenticateAndGetJWT(
		model.UserMultiAuthenticationCode(c.QueryParam("code")),
		model.UserSessionId(req.SessionId),
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.NewUserMultiAuthenticateAndGetJWTResponse(token))
}
