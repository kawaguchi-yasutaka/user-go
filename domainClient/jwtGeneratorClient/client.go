package jwtGeneratorClient

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
	"user-go/infra/jwtgenerator"
	"user-go/lib/unixtime"
)

type JwtGeneratorClient struct {
	jwtGenerator jwtgenerator.JwtGenerator
}

var _ interfaces.IJwtGeneratorClient = JwtGeneratorClient{}

func NewJwtGeneratorClient(generator jwtgenerator.JwtGenerator) JwtGeneratorClient {
	return JwtGeneratorClient{
		jwtGenerator: generator,
	}
}

func (c JwtGeneratorClient) GenerateToken(userId model.UserID, exp unixtime.UnixTime) (string, error) {
	return c.jwtGenerator.Generate(userId, exp)
}
