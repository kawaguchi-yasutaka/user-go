package model

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	mathRand "math/rand"
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

func NewMultiAuthenticationCode() (UserMultiAuthenticationCode, UserMultiAuthenticationCodeExpiresAt, error) {
	b := make([]byte, 64)
	if _, err := mathRand.Read(b); err != nil {
		return UserMultiAuthenticationCode(""), UserMultiAuthenticationCodeExpiresAt(0), err
	}
	return UserMultiAuthenticationCode(
			base64.URLEncoding.EncodeToString(b)),
		UserMultiAuthenticationCodeExpiresAt(unixtime.NewUnixTime(time.Now().Add(time.Duration(24) * time.Hour))),
		nil
}

func NewUserSessionId() (UserSessionId, UserSessionIdExpiresAt, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return UserSessionId(""), UserSessionIdExpiresAt(0), err
	}
	return UserSessionId(
			base64.URLEncoding.EncodeToString(b)),
		UserSessionIdExpiresAt(unixtime.NewUnixTime(time.Now().Add(time.Duration(24) * time.Hour))),
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

func (remember UserRemember) ValidateMultiAuthenticationCode(code UserMultiAuthenticationCode) error {
	if code != remember.MultiAuthenticationCode {
		return InvalidUserMultiAuthenticationCode("invalid code")
	}
	if unixtime.UnixTime(remember.MultiAuthenticationCodeExpiresAt) <= unixtime.Now() {
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
