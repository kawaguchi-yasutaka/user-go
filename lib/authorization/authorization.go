package authorization

import (
	"net/http"
	"user-go/lib/myerror"
)

type TokenString string

type Authorization struct {
	UserID      int64
	TokenString TokenString
}

const (
	ErrorTokenStringExpired myerror.ErrorType = "token_string_expired"
	ErrorInvalidTokenString myerror.ErrorType = "invalid_token_string"
)

func TokenStringExpired() myerror.CustomError {
	return myerror.NewCustomError("token expired", ErrorTokenStringExpired, http.StatusBadRequest)
}

func InvalidTokenString() myerror.CustomError {
	return myerror.NewCustomError("invalid token", ErrorInvalidTokenString, http.StatusBadRequest)
}

func NewAuthorization(usreId int64, tokenString TokenString) Authorization {
	return Authorization{
		UserID:      usreId,
		TokenString: tokenString,
	}
}
