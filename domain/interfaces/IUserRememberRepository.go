package interfaces

import "user-go/domain/model"

type IUserRemenberRepository interface {
	Save(userRemember model.UserRemember) error
	FindBySessionId(sessionId model.UserSessionId) (model.UserRemember, error)
}
