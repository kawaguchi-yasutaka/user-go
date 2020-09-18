package mysql

import (
	"gorm.io/gorm"
	"time"
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type userRepository struct {
	db *gorm.DB
}

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"not null;default: ''"`
	Status    string    `gorm:"not null;default: ''"`
	CreatedAt time.Time `gorm:"not null;"`
	UpdatedAt time.Time `gorm:"not null"`
}

func FromUserModel(user model.User) User {
	return User{
		ID:     int64(user.ID),
		Email:  string(user.Email),
		Status: string(user.Status),
	}
}

func NewUserRepository(db *gorm.DB) interfaces.IUserRepository {
	return userRepository{
		db: db,
	}
}

func (repo userRepository) Create(user model.User, password model.UserPasswordDigest) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		u := FromUserModel(user)
		if err := tx.Create(&u).Error; err != nil {
			return err
		}
		uPassword := FromUserPasswordModel(model.NewUserPassowrd(model.UserID(u.ID), password))
		if err := tx.Create(&uPassword).Error; err != nil {
			return err
		}
		return nil
	})
}
