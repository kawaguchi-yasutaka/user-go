package interfaces

import (
	"user-go/domain/model"
	"user-go/lib/unixtime"
)

type IJwtGeneratorClient interface {
	GenerateToken(userId model.UserID, exp unixtime.UnixTime) (string, error)
}
