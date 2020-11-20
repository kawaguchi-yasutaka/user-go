package service

import (
	"encoding/base64"
	"reflect"
	"testing"
	"time"
	"user-go/domain/model"
	"user-go/domainClient/randGenerator"
	"user-go/domainClient/userMailer"
	"user-go/infra/hasher"
	"user-go/infra/mailer"
	"user-go/infra/mysql"
	"user-go/infra/timekeeper"
	"user-go/lib/myerror"
	"user-go/lib/unixtime"
)

func initUserService() UserService {
	return UserService{
		userRepository:               mysql.NewUserRepositoryMock(),
		userAuthenticationRepository: mysql.UserAuthenticationRepositoryMock{},
		userRememberRepository:       mysql.NewUserRememberRepositoryMock(),
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

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		caseName         string
		email            model.UserEmail
		password         model.UserRawPassword
		want             model.UserSessionId
		wantErr          error
		wantUserRemember *model.UserRemember
		mailSended       bool
	}{
		{
			caseName: "ログイン成功",
			email:    "1@test.com",
			password: "password",
			want: model.UserSessionId(
				base64.URLEncoding.EncodeToString(
					[]byte("randbyte"),
				),
			),
			wantErr: nil,
			wantUserRemember: &model.UserRemember{
				UserID: 1,
				SessionId: model.UserSessionId(
					base64.URLEncoding.EncodeToString(
						[]byte("randbyte"),
					),
				),
				SessionIdExpiresAt: model.UserSessionIdExpiresAt(
					unixtime.UnixTime(time.Hour * 24),
				),
				MultiAuthenticationCode: model.UserMultiAuthenticationCode(
					base64.URLEncoding.EncodeToString(
						[]byte("randbyte"),
					),
				),
				MultiAuthenticationCodeExpiresAt: model.UserMultiAuthenticationCodeExpiresAt(
					unixtime.UnixTime(time.Hour * 24),
				),
				AuthenticationState: model.UserAuthenticationStatePendding,
			},
			mailSended: true,
		},
		{
			caseName:         "ユーザー見つからない",
			email:            "2@test.com",
			password:         "password",
			want:             model.UserSessionId(""),
			wantErr:          model.UserNotFound(),
			wantUserRemember: nil,
			mailSended:       false,
		},
		{
			caseName:         "パスワードが違う",
			email:            "3@test.com",
			password:         "not correct",
			want:             model.UserSessionId(""),
			wantErr:          model.IncorrectUserPassword(""),
			wantUserRemember: nil,
			mailSended:       false,
		},
	}

	for _, tt := range tests {
		service := initUserService()
		users := map[model.UserID]model.User{1: {ID: 1, Email: "1@test.com"}, 3: {ID: 3, Email: "3@test.com"}}
		userRemembers := map[model.UserSessionId]model.UserRemember{}
		userAuthentications := map[model.UserID]model.UserAuthentication{
			1: {UserID: 1, PasswordDigest: "password"},
			3: {UserID: 3, PasswordDigest: "password"},
		}

		multiAuthenticationCodeSendedCount := 0

		service.userRepository = mysql.UserRepositoryMock{Users: users}
		service.userRememberRepository = mysql.UserRememberRepositoryMock{UserRemembers: userRemembers}
		service.randGenerator = randGenerator.RandGeneratorMock{RandByte: []byte("randbyte")}
		service.userMailer = userMailer.UserMailerMock{MultiAuthenticationCodeSendedCount: &multiAuthenticationCodeSendedCount}
		service.userAuthenticationRepository = mysql.UserAuthenticationRepositoryMock{UserAuthentications: userAuthentications}

		result, err := service.Login(tt.email, tt.password)
		if !myerror.EqualErrorType(err, tt.wantErr) {
			t.Errorf("casename: %v, err: %v,wantErr: %v", tt.caseName, err, tt.wantErr)
		}
		if result != tt.want {
			t.Errorf("casename: %v, result: %v,want: %v", tt.caseName, result, tt.want)
		}
		if tt.wantUserRemember != nil {
			if userRemembers[result] != *tt.wantUserRemember {
				t.Errorf("casename: %v, user remember: %v,want: %v", tt.caseName, result, tt.wantUserRemember)
			}
		}
		if tt.mailSended {
			if multiAuthenticationCodeSendedCount != 1 {
				t.Errorf("casename: %v, mail send failed", tt.caseName)
			}
		}
	}
}

func TestUserService_MultiAuthenticate(t *testing.T) {
	tests := []struct {
		caseName  string
		code      model.UserMultiAuthenticationCode
		sessionID model.UserSessionId
		wantErr   error
	}{
		{
			caseName:  "2段階認証成功",
			code:      model.UserMultiAuthenticationCode("1"),
			sessionID: model.UserSessionId("1"),
			wantErr:   nil,
		},
		{
			caseName:  "セッションIDに紐づいてる二段階認証コードとリクエストのコードが不一致",
			code:      model.UserMultiAuthenticationCode("3"),
			sessionID: model.UserSessionId("2"),
			wantErr:   model.InvalidUserMultiAuthenticationCode(""),
		},
		{
			caseName:  "二段階認証コードが、期限切れ",
			code:      model.UserMultiAuthenticationCode("3"),
			sessionID: model.UserSessionId("3"),
			wantErr:   model.ExpiredUserMultiAuthenticationCode(),
		},
		{
			caseName:  "既に二段階認証済み",
			code:      model.UserMultiAuthenticationCode("4"),
			sessionID: model.UserSessionId("4"),
			wantErr:   model.AlreadyMultiAuthenticated(""),
		},
	}

	for _, tt := range tests {
		service := initUserService()
		userRemembers := map[model.UserSessionId]model.UserRemember{
			model.UserSessionId("1"): {SessionId: "1", MultiAuthenticationCode: "1", MultiAuthenticationCodeExpiresAt: 1},
			model.UserSessionId("2"): {SessionId: "2", MultiAuthenticationCode: "2", MultiAuthenticationCodeExpiresAt: 1},
			model.UserSessionId("3"): {SessionId: "3", MultiAuthenticationCode: "3", MultiAuthenticationCodeExpiresAt: 0},
			model.UserSessionId("4"): {SessionId: "4", MultiAuthenticationCode: "4", MultiAuthenticationCodeExpiresAt: 1, AuthenticationState: model.UserAuthenticationStateComplete},
		}
		service.userRememberRepository = mysql.UserRememberRepositoryMock{
			UserRemembers: userRemembers,
		}
		service.timekeeper = timekeeper.TimeKeeperMock{N: 0}

		if err := service.MultiAuthenticate(tt.code, tt.sessionID); !myerror.EqualErrorType(err, tt.wantErr) {
			t.Errorf("casename: %v, err: %v,wantErr: %v", tt.caseName, err, tt.wantErr)
		}
		if tt.wantErr == nil {
			v, ok := userRemembers[tt.sessionID]
			if !ok {
				t.Errorf("casename: %v, user remember not found", tt.caseName)
			}
			if !v.IsComplete() {
				t.Errorf("casename: %v, multi authenticate not complete %v", tt.caseName, v.AuthenticationState)
			}
		}
	}
}

func TestUserService_GenerateJwtToken(t *testing.T) {

}
