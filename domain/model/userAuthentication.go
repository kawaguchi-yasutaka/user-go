package model

type UserAuthentication struct {
	UserID         UserID
	PasswordDigest UserPasswordDigest
}

type (
	UserRawPassword    string
	UserPasswordDigest string
)

func NewUserAuthentication(userId UserID, passwordDigest UserPasswordDigest) UserAuthentication {
	return UserAuthentication{
		UserID:         userId,
		PasswordDigest: passwordDigest,
	}
}

func NewUserRawPassword(password string) (UserRawPassword, error) {
	if err := Validate.Var(password, "required,min=8"); err != nil {
		return "", err
	}
	return UserRawPassword(password), nil
}
