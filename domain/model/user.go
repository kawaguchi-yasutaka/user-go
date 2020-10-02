package model

import "fmt"

type User struct {
	ID     UserID
	Email  UserEmail
	Status UserStatus
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

func NewUser(email UserEmail) User {
	return User{
		Email:  email,
		Status: UserStatusInitialized,
	}
}

func NewUserEmail(email string) (UserEmail, error) {
	if err := Validate.Var(email, "required,email"); err != nil {
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
}
