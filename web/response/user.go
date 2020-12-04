package response

type UserMultiAuthenticateAndGetJWTResponse struct {
	Token string `json:"token"`
}

func NewUserMultiAuthenticateAndGetJWTResponse(token string) UserMultiAuthenticateAndGetJWTResponse {
	return UserMultiAuthenticateAndGetJWTResponse{
		Token: token,
	}
}
