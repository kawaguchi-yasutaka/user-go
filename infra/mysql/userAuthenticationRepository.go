package mysql

import (
	"errors"
	"gorm.io/gorm"
	"user-go/domain/interfaces"
	"user-go/domain/model"
	"user-go/lib/myerror"
)

type UserAuthentication struct {
	UserID                 int64   `gorm:"primaryKey"`
	PasswordDigest         string  `gorm:"not null;default:''"`
	ActivationCode         string  `gorm:"not null;default:''"`
	ActivationCodeExpireAt int64   `gorm:"not null;default:0"`
	SessionId              *string `gorm:"unique"`
	SessionIdExpiresAt     int64   `gorm:"not null;default:0"`
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
	var sessionId *string
	if auth.SessionId != nil {
		val := string(*auth.SessionId)
		sessionId = &val
	}
	return UserAuthentication{
		UserID:                 int64(auth.UserID),
		PasswordDigest:         string(auth.PasswordDigest),
		ActivationCode:         string(auth.ActivationCode),
		ActivationCodeExpireAt: int64(auth.ActivationCodeExpiresAt),
		SessionId:              sessionId,
		SessionIdExpiresAt:     int64(auth.SessionIdExpiresAt),
	}
}

func (authentication UserAuthentication) ToModel() model.UserAuthentication {
	var sessionId *model.UserSessionId
	if authentication.SessionId != nil {
		val := model.UserSessionId(*authentication.SessionId)
		sessionId = &val
	}
	return model.UserAuthentication{
		UserID:                  model.UserID(authentication.UserID),
		PasswordDigest:          model.UserPasswordDigest(authentication.PasswordDigest),
		ActivationCode:          model.UserActivationCode(authentication.ActivationCode),
		ActivationCodeExpiresAt: model.UserActivationCodeExpiresAt(authentication.ActivationCodeExpireAt),
		SessionId:               sessionId,
		SessionIdExpiresAt:      model.UserSessionIdExpiresAt(authentication.SessionIdExpiresAt),
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

func (repo UserAuthenticationRepository) FindByActivateCode(
	code model.UserActivationCode,
) (model.UserAuthentication, error) {
	auth := UserAuthentication{}
	if result := repo.db.Where("activation_code = ?", string(code)).First(&auth); result.Error != nil {
		return model.UserAuthentication{}, result.Error
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
