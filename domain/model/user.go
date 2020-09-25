package model

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
