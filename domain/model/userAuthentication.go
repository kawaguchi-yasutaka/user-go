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
