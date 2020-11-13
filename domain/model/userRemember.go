package model

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
	"user-go/lib/myerror"
	"user-go/lib/unixtime"
)

//Todo 名前two factor authenticationに変える。
type UserRemember struct {
	SessionId                        UserSessionId
	UserID                           UserID
	MultiAuthenticationCode          UserMultiAuthenticationCode
	MultiAuthenticationCodeExpiresAt UserMultiAuthenticationCodeExpiresAt
	SessionIdExpiresAt               UserSessionIdExpiresAt
	AuthenticationState              UserAuthenticationState
}

type (
	UserAuthenticationState              string
	UserMultiAuthenticationCode          string
	UserMultiAuthenticationCodeExpiresAt unixtime.UnixTime
	UserSessionId                        string
	UserSessionIdExpiresAt               unixtime.UnixTime
)

const (
	UserAuthenticationStatePendding UserAuthenticationState = "pending"
	UserAuthenticationStateComplete UserAuthenticationState = "complete"
)

const (
	ErrorExpiredUserMultiAuthenticationCode myerror.ErrorType = "expired_user_multi_authentication_code"
	ErrorInvalidUserMultiAuthenticationCode myerror.ErrorType = "invalid_user_multi_authentication_code"
	ErrorNotCompleteUserAuthentication      myerror.ErrorType = "not_complete_user_authentication"
	ErrorUserRememberNotFound               myerror.ErrorType = "user_remember_not_found"
	ErrorAlreadyMultiAuthenticated          myerror.ErrorType = "already_multi_authenticated"
)

func ExpiredUserMultiAuthenticationCode() myerror.CustomError {
	return myerror.NewCustomError("expired multi authentication code", ErrorExpiredUserMultiAuthenticationCode, http.StatusBadRequest)
}

func InvalidUserMultiAuthenticationCode(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorInvalidUserMultiAuthenticationCode, http.StatusBadRequest)
}

func NotCompleteUserAuthentication(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorNotCompleteUserAuthentication, http.StatusUnauthorized)
}

func UserRememberNotFound(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorUserRememberNotFound, http.StatusNotFound)
}

func AlreadyMultiAuthenticated(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorAlreadyMultiAuthenticated, http.StatusBadRequest)
}

func NewMultiAuthenticationCode(rand []byte, now unixtime.UnixTime) (UserMultiAuthenticationCode, UserMultiAuthenticationCodeExpiresAt, error) {
	return UserMultiAuthenticationCode(base64.URLEncoding.EncodeToString(rand)),
		UserMultiAuthenticationCodeExpiresAt(now + unixtime.UnixTime(time.Duration(24)*time.Hour)),
		nil
}

func NewUserSessionId(rand []byte, now unixtime.UnixTime) (UserSessionId, UserSessionIdExpiresAt, error) {
	return UserSessionId(
			base64.URLEncoding.EncodeToString(rand)),
		UserSessionIdExpiresAt(now + unixtime.UnixTime(time.Duration(24)*time.Hour)),
		nil
}

func (remember *UserRemember) UpdateMultiAuthenticationInfo(code UserMultiAuthenticationCode, expiresAt UserMultiAuthenticationCodeExpiresAt) {
	remember.MultiAuthenticationCode = code
	remember.MultiAuthenticationCodeExpiresAt = expiresAt
}

func (remember *UserRemember) UpdateSessionInfo(id UserSessionId, expiresAt UserSessionIdExpiresAt) {
	remember.SessionId = id
	remember.SessionIdExpiresAt = expiresAt
}

func (remember UserRemember) ValidateMultiAuthenticationCode(code UserMultiAuthenticationCode, now unixtime.UnixTime) error {
	if code != remember.MultiAuthenticationCode {
		return InvalidUserMultiAuthenticationCode("invalid code")
	}
	if unixtime.UnixTime(remember.MultiAuthenticationCodeExpiresAt) <= now {
		return ExpiredUserMultiAuthenticationCode()
	}
	return nil
}

func (remember UserRemember) ValidateSession() error {
	if remember.AuthenticationState == UserAuthenticationStatePendding {
		return NotCompleteUserAuthentication(fmt.Sprintf("user id %v is authoraication not complete", remember.UserID))
	}
	if unixtime.UnixTime(remember.SessionIdExpiresAt) <= unixtime.Now() {
		return UserUnauthorized("session id is expired")
	}
	return nil
}

func NewUserRememberBySingleFactorAuthentication(
	userId UserID,
	sessionId UserSessionId,
	sessionIdExpiresAt UserSessionIdExpiresAt,
	code UserMultiAuthenticationCode,
	codeExpiresAt UserMultiAuthenticationCodeExpiresAt,
) UserRemember {
	return UserRemember{
		UserID:                           userId,
		SessionId:                        sessionId,
		SessionIdExpiresAt:               sessionIdExpiresAt,
		MultiAuthenticationCode:          code,
		MultiAuthenticationCodeExpiresAt: codeExpiresAt,
		AuthenticationState:              UserAuthenticationStatePendding,
	}
}

func (remember *UserRemember) Completed() {
	remember.AuthenticationState = UserAuthenticationStateComplete
}

func (remember UserRemember) IsComplete() bool {
	return remember.AuthenticationState == UserAuthenticationStateComplete
}
