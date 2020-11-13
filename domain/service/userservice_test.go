package service

import (
	"encoding/base64"
	"reflect"
	"testing"
	"time"
	"user-go/domain/model"
	"user-go/domainClient/randGenerator"
	"user-go/infra/hasher"
	"user-go/infra/mailer"
	"user-go/infra/mysql"
	"user-go/infra/timekeeper"
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
		randGenerator:                randGenerator.NewRandGeneratorMock(),
		timekeeper:                   timekeeper.NewTimeKeeperMock(),
	}
}

func TestUserService_Create(t *testing.T) {
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
				ActivationCode:          model.UserActivationCode(base64.URLEncoding.EncodeToString([]byte("activate code"))),
				ActivationCodeExpiresAt: model.UserActivationCodeExpiresAt(unixtime.UnixTime(time.Hour * 24)),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		service := initUserService()
		users := make(map[model.UserID]model.User)
		userAuthentications := make(map[model.UserID]model.UserAuthentication)
		service.userRepository = mysql.UserRepositoryMock{Users: users, UserAuthentications: userAuthentications}
		service.userAuthenticationRepository = mysql.UserAuthenticationRepositoryMock{UserAuthentications: userAuthentications}

		service.timekeeper = timekeeper.TimeKeeperMock{N: 0}
		service.randGenerator = randGenerator.RandGeneratorMock{RandByte: []byte("activate code")}

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
			id:         model.UserID(2),
			wantStatus: model.UserStatusActivated,
			wantErr:    model.AlreadyActivated(""),
		},
		{
			caseName:   "アクティベートコードが期限切れ",
			code:       model.UserActivationCode("3"),
			id:         model.UserID(3),
			wantStatus: model.UserStatusActivated,
			wantErr:    model.ExpiredUserActivationCode(),
		},
	}

	for _, tt := range tests {
		service := initUserService()
		users := map[model.UserID]model.User{1: {ID: 1}, 2: {ID: 2, Status: model.UserStatusActivated}, 3: {ID: 3}}
		userAuthentications := map[model.UserID]model.UserAuthentication{
			1: {UserID: 1, ActivationCode: "1", ActivationCodeExpiresAt: 1},
			2: {UserID: 2, ActivationCode: "2", ActivationCodeExpiresAt: 1},
			3: {UserID: 3, ActivationCode: "3", ActivationCodeExpiresAt: 0},
		}
		service.userRepository = mysql.UserRepositoryMock{Users: users, UserAuthentications: userAuthentications}
		service.userAuthenticationRepository = mysql.UserAuthenticationRepositoryMock{UserAuthentications: userAuthentications}
		service.timekeeper = timekeeper.TimeKeeperMock{N: 0}

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
