package interfaces

import "user-go/lib/authorization"

type IJwtHandlerClient interface {
	Parse(token authorization.TokenString) (authorization.Authorization, error)
}
