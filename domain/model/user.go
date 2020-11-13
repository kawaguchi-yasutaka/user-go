package model

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
	"user-go/lib/myerror"
)

type User struct {
	ID        UserID
	Email     UserEmail
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type (
	UserID     int64
	UserStatus string
	UserEmail  string
)

const (
	UserStatusInitialized UserStatus = "initialized"
	UserStatusActivated   UserStatus = "activated" //認証済み
)

var UserStates = []UserStatus{UserStatusInitialized, UserStatusActivated}

const (
	ErrorUserNotFound           myerror.ErrorType = "user_not_found"
	ErrorRequiredUserEmail      myerror.ErrorType = "required_user_email"
	ErrorInvalidFormatUserEmail myerror.ErrorType = "invalid_format_user_email"
	ErrorIncorrectUserPassword  myerror.ErrorType = "incorrect_user_password"
	ErrorAlreadyActivated       myerror.ErrorType = "already_activated"
)

func NewUser(email UserEmail) User {
	return User{
		Email:  email,
		Status: UserStatusInitialized,
	}
}
func UserNotFound() myerror.CustomError {
	return myerror.NewCustomError("user not found", ErrorUserNotFound, http.StatusNotFound)
}

func RequiredUserEmail(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorRequiredUserEmail, http.StatusBadRequest)
}

func InvalidFormatUserEmail(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorInvalidFormatUserEmail, http.StatusBadRequest)
}

func IncorrectUserPassword(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorIncorrectUserPassword, http.StatusBadRequest)
}

func AlreadyActivated(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorAlreadyActivated, http.StatusBadRequest)
}

func NewUserEmail(email string) (UserEmail, error) {
	if err := Validate.Var(email, "required,email"); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			switch fieldName {
			case "required":
				return "", RequiredUserEmail("")
			case "email":
				return "", InvalidFormatUserEmail(fmt.Sprintf(
					"%s is invalid email format",
					email,
				))
			}
		}
		return "", err
	}
	return UserEmail(email), nil
}

func NewUserStatus(status string) (UserStatus, error) {
	for _, v := range UserStates {
		if string(v) == status {
			return v, nil
		}
	}
	return "", fmt.Errorf("invalid status %v", status)
}

func (user *User) IsActivated() bool {
	return user.Status == UserStatusActivated
}

func (user *User) Activate() {
	user.Status = UserStatusActivated
	user.UpdatedAt = time.Now()
}
