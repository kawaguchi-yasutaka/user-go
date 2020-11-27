package interfaces

import "user-go/domain/model"

type IJwtGeneratorClient interface {
	GenerateToken(userId model.UserID, exp string) (string, error)
}
