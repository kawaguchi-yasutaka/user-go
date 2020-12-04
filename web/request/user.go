package request

type UserCreateRequest struct {
	Email    string `json:email;`
	Password string `json:password;`
}

type UserLoginRequest struct {
	Email    string `json:email;`
	Password string `json:password;`
}

type UserMultiAuthenticateRequest struct {
	SessionId string `json:sessionId;`
}

type UserMultiAuthenticateAndGetJWTRequest struct {
	SessionId string `json:sessionId;`
}
