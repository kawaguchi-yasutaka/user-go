package model

type UserPassword struct {
	UserID         UserID
	PasswordDigest UserPasswordDigest
}

type (
	UserRawPassword    string
	UserPasswordDigest string
)

func NewUserPassowrd(userId UserID, passwordDigest UserPasswordDigest) UserPassword {
	return UserPassword{
		UserID:         userId,
		PasswordDigest: passwordDigest,
	}
}
