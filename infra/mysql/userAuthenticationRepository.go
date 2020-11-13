package mysql

import (
	"errors"
	"gorm.io/gorm"
	"user-go/domain/interfaces"
	"user-go/domain/model"
	"user-go/lib/myerror"
)

type UserAuthentication struct {
	UserID                 int64  `gorm:"primaryKey"`
	PasswordDigest         string `gorm:"not null;default:''"`
	ActivationCode         string `gorm:"not null;default:''"`
	ActivationCodeExpireAt int64  `gorm:"not null;default:0"`
}

type UserAuthenticationRepository struct {
	db *gorm.DB
}

func NewUserAuthenticationRepository(db *gorm.DB) interfaces.IUserAuthenticationRepository {
	return UserAuthenticationRepository{
		db: db,
	}
}

func FromUserAuthenticationModel(auth model.UserAuthentication) UserAuthentication {
	return UserAuthentication{
		UserID:                 int64(auth.UserID),
		PasswordDigest:         string(auth.PasswordDigest),
		ActivationCode:         string(auth.ActivationCode),
		ActivationCodeExpireAt: int64(auth.ActivationCodeExpiresAt),
	}
}

func (authentication UserAuthentication) ToModel() model.UserAuthentication {
	return model.UserAuthentication{
		UserID:                  model.UserID(authentication.UserID),
		PasswordDigest:          model.UserPasswordDigest(authentication.PasswordDigest),
		ActivationCode:          model.UserActivationCode(authentication.ActivationCode),
		ActivationCodeExpiresAt: model.UserActivationCodeExpiresAt(authentication.ActivationCodeExpireAt),
	}
}

func (repo UserAuthenticationRepository) Save(authentication model.UserAuthentication) error {
	a := FromUserAuthenticationModel(authentication)
	if err := repo.db.Save(&a).Debug().Error; err != nil {
		return myerror.DBError(err)
	}
	return nil
}

func (repo UserAuthenticationRepository) FindByUserID(
	userID model.UserID,
) (model.UserAuthentication, error) {
	auth := UserAuthentication{}
	result := repo.db.Where("user_id = ?", int64(userID)).First(&auth)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.UserAuthentication{}, model.UserAuthenticationNotFound()
	}
	if result.Error != nil {
		return model.UserAuthentication{}, myerror.DBError(result.Error)
	}
	return auth.ToModel(), nil
}

func (repo UserAuthenticationRepository) FindByActivateCodeAndUserID(
	code model.UserActivationCode,
	id model.UserID,
) (model.UserAuthentication, error) {
	auth := UserAuthentication{}
	result := repo.db.
		Where("activation_code = ? AND user_id = ?", string(code), int64(id)).
		First(&auth)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.UserAuthentication{}, model.UserAuthenticationNotFound()
	}
	if result.Error != nil {
		return model.UserAuthentication{}, myerror.DBError(result.Error)
	}
	return auth.ToModel(), nil
}

func (repo UserAuthenticationRepository) FindBySessionId(
	sessionId model.UserSessionId,
) (model.UserAuthentication, error) {
	auth := UserAuthentication{}
	result := repo.db.Where("session_id = ?", string(sessionId)).First(&auth)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.UserAuthentication{}, model.UserAuthenticationNotFound()
	}
	if result.Error != nil {
		return model.UserAuthentication{}, myerror.DBError(result.Error)
	}
	return auth.ToModel(), nil
}

func (repo UserAuthenticationRepository) FindByMultiAuthenticateCode(
	code model.UserMultiAuthenticationCode,
	id model.UserID,
) (model.UserAuthentication, error) {
	auth := UserAuthentication{}
	result := repo.db.Where("multi_authentication_code = ? AND user_id = ?", string(code), int64(id)).First(&auth)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.UserAuthentication{}, model.UserAuthenticationNotFound()
	}
	if result.Error != nil {
		return model.UserAuthentication{}, myerror.DBError(result.Error)
	}
	return auth.ToModel(), nil
}
