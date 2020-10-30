package service

import (
	"reflect"
	"testing"
	"user-go/domain/model"
	"user-go/infra/hasher"
	"user-go/infra/mailer"
	"user-go/infra/mysql"
	"user-go/lib/myerror"
	"user-go/lib/unixtime"
)

func initUserService() UserService {
	return UserService{
		userRepository:               mysql.UserRepositoryMock{},
		userAuthenticationRepository: mysql.UserAuthenticationRepositoryMock{},
		userRememberRepository:       mysql.UserRememberRepository{},
		hasher:                       hasher.HasherMock{},
		userMailer:                   mailer.MailerMock{},
	}
}

func TestUserService_Create(t *testing.T) {
	service := initUserService()
	users := make(map[model.UserID]model.User)
	userAuthentications := make(map[model.UserID]model.UserAuthentication)
	service.userRepository = mysql.UserRepositoryMock{Users: users, UserAuthentications: userAuthentications}
	service.userAuthenticationRepository = mysql.UserAuthenticationRepositoryMock{UserAuthentications: userAuthentications}

	tests := []struct {
		caseName               string
		email                  model.UserEmail
		password               model.UserRawPassword
		wantUser               model.User
		wantUserAuthentication model.UserAuthentication
		wantErr                error
	}{
		{
			caseName: "新規作成",
			email:    model.UserEmail("test@test.com"),
			password: model.UserRawPassword("password"),
			wantUser: model.User{
				ID:     1,
				Email:  model.UserEmail("test@test.com"),
				Status: model.UserStatusInitialized,
			},
			wantUserAuthentication: model.UserAuthentication{
				UserID:                  1,
				PasswordDigest:          "password",
				ActivationCode:          model.UserActivationCode(""),
				ActivationCodeExpiresAt: model.UserActivationCodeExpiresAt(unixtime.UnixTime(0)),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		userID, err := service.Create(tt.email, tt.password)
		if !myerror.EqualErrorType(err, tt.wantErr) {
			t.Errorf("casename: %v, err: %v,wantErr: %v", tt.caseName, err, tt.wantErr)
		}
		if userID != 0 {
			u, ok := users[userID]
			if !ok {
				t.Errorf("casename: %v, userId %v user not created", tt.caseName, userID)
			}
			if !reflect.DeepEqual(u, tt.wantUser) {
				t.Errorf("casename: %v, user %v want user %v", tt.caseName, u, tt.wantUser)
			}
			a, ok := userAuthentications[userID]
			if !ok {
				t.Errorf("casename: %v, userId %v userauthentication not created", tt.caseName, userID)
			}
			if !reflect.DeepEqual(a, tt.wantUserAuthentication) {
				t.Errorf("casename: %v, user %v want user %v", tt.caseName, a, tt.wantUserAuthentication)
			}
		}
	}
}

func TestUserService_Activate(t *testing.T) {
	service := initUserService()
	users := map[model.UserID]model.User{1: {ID: 1}, 2: {ID: 2, Status: model.UserStatusActivated}, 3: {ID: 3}}
	userAuthentications := map[model.UserID]model.UserAuthentication{
		1: {UserID: 1, ActivationCode: "1"},
		2: {UserID: 2, ActivationCode: "2"},
		3: {UserID: 3, ActivationCode: "3"},
	}
	service.userRepository = mysql.UserRepositoryMock{Users: users, UserAuthentications: userAuthentications}
	service.userAuthenticationRepository = mysql.UserAuthenticationRepositoryMock{UserAuthentications: userAuthentications}

	tests := []struct {
		caseName   string
		code       model.UserActivationCode
		id         model.UserID
		wantStatus model.UserStatus
		wantErr    error
	}{
		{
			caseName:   "有効化成功",
			code:       model.UserActivationCode("1"),
			id:         model.UserID(1),
			wantStatus: model.UserStatusActivated,
			wantErr:    nil,
		},
		{
			caseName:   "既に有効化済み",
			code:       model.UserActivationCode("2"),
			id:         model.UserID(1),
			wantStatus: model.UserStatusActivated,
			wantErr:    nil,
		},
		{
			caseName:   "アクティベートコードが期限切れ",
			code:       model.UserActivationCode("3"),
			id:         model.UserID(1),
			wantStatus: model.UserStatusActivated,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		err := service.Activate(tt.code, tt.id)
		if !myerror.EqualErrorType(err, tt.wantErr) {
			t.Errorf("casename: %v, err: %v,wantErr: %v", tt.caseName, err, tt.wantErr)
		}
		if err == nil {
			user, ok := users[tt.id]
			if !ok {
				t.Errorf("casename: %v, userId %v not found", tt.caseName, tt.id)
			}
			if user.Status != tt.wantStatus {
				t.Errorf("casename: %v, status %v want status %v", tt.caseName, user.Status, tt.wantStatus)
			}
		}
	}
}
