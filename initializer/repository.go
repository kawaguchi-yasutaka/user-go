package initializer

import (
	"gorm.io/gorm"
	"user-go/domain/interfaces"
	"user-go/infra/mysql"
)

type Repository struct {
	userRepository interfaces.IUserRepository
}

func InitRepository(db *gorm.DB) Repository {
	return Repository{
		userRepository: mysql.NewUserRepository(db),
	}
}
