package mysql

import (
	"errors"
	"gorm.io/gorm"
	"user-go/domain/interfaces"
	"user-go/domain/model"
	"user-go/lib/myerror"
)

type UserRemember struct {
	SessionId                        string `gorm:"primaryKey"`
	SessionIdExpiresAt               int64  `gorm:"not null;default:0"`
	UserID                           int64  `gorm:"not null;default:0"`
	MultiAuthenticationCode          string `gorm:"not null;default:''"`
	MultiAuthenticationCodeExpiresAt int64  `gorm:"not null;default:0"`
	AuthenticationState              string `gorm:"not null;default:''"`
}

type UserRememberRepository struct {
	db *gorm.DB
}

func (remember UserRemember) ToModel() model.UserRemember {
	return model.UserRemember{
		SessionId:                        model.UserSessionId(remember.SessionId),
		SessionIdExpiresAt:               model.UserSessionIdExpiresAt(remember.SessionIdExpiresAt),
		UserID:                           model.UserID(remember.UserID),
		MultiAuthenticationCode:          model.UserMultiAuthenticationCode(remember.MultiAuthenticationCode),
		MultiAuthenticationCodeExpiresAt: model.UserMultiAuthenticationCodeExpiresAt(remember.MultiAuthenticationCodeExpiresAt),
		AuthenticationState:              model.UserAuthenticationState(remember.AuthenticationState),
	}
}

func FromUserRememberModel(remember model.UserRemember) UserRemember {
	return UserRemember{
		SessionId:                        string(remember.SessionId),
		SessionIdExpiresAt:               int64(remember.SessionIdExpiresAt),
		UserID:                           int64(remember.UserID),
		MultiAuthenticationCode:          string(remember.MultiAuthenticationCode),
		MultiAuthenticationCodeExpiresAt: int64(remember.MultiAuthenticationCodeExpiresAt),
		AuthenticationState:              string(remember.AuthenticationState),
	}
}

func NewUserRememberRepository(db *gorm.DB) interfaces.IUserRemenberRepository {
	return UserRememberRepository{
		db: db,
	}
}

func (repo UserRememberRepository) Save(userRemember model.UserRemember) error {
	r := FromUserRememberModel(userRemember)
	if err := repo.db.Save(&r).Error; err != nil {
		return myerror.DBError(err)
	}
	return nil
}

func (repo UserRememberRepository) FindBySessionId(
	sessionId model.UserSessionId,
) (model.UserRemember, error) {
	remember := UserRemember{}
	result := repo.db.Where("session_id = ?", string(sessionId)).First(&remember)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.UserRemember{}, model.UserAuthenticationNotFound()
	}
	if result.Error != nil {
		return model.UserRemember{}, myerror.DBError(result.Error)
	}
	return remember.ToModel(), nil
}
