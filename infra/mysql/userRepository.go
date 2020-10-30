package mysql

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"user-go/domain/interfaces"
	"user-go/domain/model"
	"user-go/lib/myerror"
)

type userRepository struct {
	db *gorm.DB
}

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"not null;unique;default:'';index:user_email_idx"`
	Status    string    `gorm:"not null;default:''"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func FromUserModel(user model.User) User {
	return User{
		ID:        int64(user.ID),
		Email:     string(user.Email),
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
func (user User) ToModel() (model.User, error) {
	email, err := model.NewUserEmail(user.Email)
	if err != nil {
		return model.User{}, err
	}
	status, err := model.NewUserStatus(user.Status)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		ID:        model.UserID(user.ID),
		Email:     email,
		Status:    status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func NewUserRepository(db *gorm.DB) interfaces.IUserRepository {
	return userRepository{
		db: db,
	}
}

func (repo userRepository) Create(user model.User, password model.UserPasswordDigest) (model.UserID, error) {
	u := FromUserModel(user)
	return model.UserID(u.ID), repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&u).Error; err != nil {
			return myerror.DBError(err)
		}
		uPassword := FromUserAuthenticationModel(model.NewUserAuthentication(model.UserID(u.ID), password))
		if err := tx.Create(&uPassword).Error; err != nil {
			return myerror.DBError(err)
		}
		return nil
	})
}

func (repo userRepository) Save(user model.User) error {
	u := FromUserModel(user)
	if err := repo.db.Save(&u).Error; err != nil {
		return err
	}
	return nil
}

func (repo userRepository) FindById(id model.UserID) (model.User, error) {
	user := User{}
	if err := repo.db.Where("id = ?", int64(id)).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user.ToModel()
}

func (repo userRepository) FindByEmail(email model.UserEmail) (model.User, error) {
	user := User{}
	result := repo.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{}, model.UserNotFound()
	}
	if result.Error != nil {
		return model.User{}, myerror.DBError(result.Error)
	}
	return user.ToModel()
}
