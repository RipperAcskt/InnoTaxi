package handler

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/RipperAcskt/innotaxi/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestSingUp(t *testing.T) {
	type mockBehavior func(s *mocks.MockAuthRepo, user service.UserSingUp)
	type fileds struct {
		authRepo  *mocks.MockAuthRepo
		tokenRepo *mocks.MockTokenRepo
	}
	test := []struct {
		name         string
		body         string
		user         service.UserSingUp
		mockBehavior mockBehavior
		expectedCode int
		expectedBody string
	}{
		{
			name: "new user",
			body: `{"name": "Ivan", "phone_number": "+7455456", "email": "ripper@algsdh", "password": "12345"}`,
			user: service.UserSingUp{
				Name:        "Ivan",
				PhoneNumber: "+7455456",
				Email:       "ripper@algsdh",
				Password:    "12345",
			},
			mockBehavior: func(s *mocks.MockAuthRepo, user service.UserSingUp) {
				s.EXPECT().CreateUser(context.Background(), user).Return(nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: "",
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fileds{
				authRepo:  mocks.NewMockAuthRepo(ctrl),
				tokenRepo: mocks.NewMockTokenRepo(ctrl),
			}
			authService := service.NewAuthSevice(f.authRepo, f.tokenRepo, "124jkhsdaf3425", &config.Config{})

			tt.user.Password, _ = authService.GenerateHash(tt.user.Password)
			tt.mockBehavior(f.authRepo, tt.user)

			service := service.Service{
				AuthService: authService,
			}
			logger, _ := zap.NewProduction()
			handler := New(&service, nil, logger)

			r := gin.New()
			r.POST("/users/auth/sing-up", handler.singUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/auth/sing-up", bytes.NewBufferString(tt.body))

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestSingIn(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config new failed: %v", err)
	}

	type mockBehavior func(s *mocks.MockAuthRepo, phone_number string)
	type fileds struct {
		authRepo  *mocks.MockAuthRepo
		tokenRepo *mocks.MockTokenRepo
	}
	test := []struct {
		name         string
		body         string
		user         service.UserSingIn
		mockBehavior mockBehavior
		expectedCode int
		expectedBody string
	}{
		{
			name: "correct password",
			body: `{"phone_number": "2", "password": "2"}`,
			user: service.UserSingIn{
				ID:          9,
				PhoneNumber: "2",
				Password:    "2",
			},
			mockBehavior: func(s *mocks.MockAuthRepo, phone_number string) {
				s.EXPECT().CheckUserByPhoneNumber(context.Background(), phone_number).Return(&service.UserSingIn{
					ID:          9,
					PhoneNumber: "2",
					Password:    string([]byte{49, 50, 52, 106, 107, 104, 115, 100, 97, 102, 51, 52, 50, 53, 218, 75, 146, 55, 186, 204, 205, 241, 156, 7, 96, 202, 183, 174, 196, 168, 53, 144, 16, 176}),
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name: "incorrect password",
			body: `{"phone_number": "+7455456", "password": "123456"}`,
			user: service.UserSingIn{
				ID:          1,
				PhoneNumber: "+7455456",
				Password:    "123456",
			},
			mockBehavior: func(s *mocks.MockAuthRepo, phone_number string) {
				s.EXPECT().CheckUserByPhoneNumber(context.Background(), phone_number).Return(&service.UserSingIn{}, nil)
			},
			expectedCode: http.StatusForbidden,
			expectedBody: "",
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fileds{
				authRepo:  mocks.NewMockAuthRepo(ctrl),
				tokenRepo: mocks.NewMockTokenRepo(ctrl),
			}
			authService := service.NewAuthSevice(f.authRepo, f.tokenRepo, "124jkhsdaf3425", &config.Config{})

			tt.mockBehavior(f.authRepo, tt.user.PhoneNumber)

			service := service.Service{
				AuthService: authService,
			}
			logger, _ := zap.NewProduction()
			handler := New(&service, cfg, logger)

			r := gin.New()
			r.POST("/users/auth/sing-in", handler.singIn)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/auth/sing-in", bytes.NewBufferString(tt.body))

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.NotEqual(t, tt.expectedBody, w.Body.String())
		})
	}

}
