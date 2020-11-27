package jwtGeneratorClient

import (
	"user-go/domain/model"
	"user-go/infra/jwtgenerator"
)

type JwtGeneratorClient struct {
	jwtGenerator jwtgenerator.JwtGenerator
}

func NewJwtGeneratorClient(generator jwtgenerator.JwtGenerator) JwtGeneratorClient {
	return JwtGeneratorClient{
		jwtGenerator: generator,
	}
}

func (c JwtGeneratorClient) GenerateToken(userId model.UserID, exp string) (string, error) {
	return c.jwtGenerator.Generate(
		map[string]interface{}{
			"exp":    exp,
			"userId": userId,
		},
	)
}
