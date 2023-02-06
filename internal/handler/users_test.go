package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RipperAcskt/innotaxi/internal/model"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/RipperAcskt/innotaxi/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestGetProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserRepo)

	test := []struct {
		name         string
		mockBehavior mockBehavior
		expectedCode int
		expectedBody string
	}{
		{
			name: "get user",
			mockBehavior: func(s *mocks.MockUserRepo) {
				s.EXPECT().GetUserById(context.Background(), "").Return(&model.User{
					Name:        "2",
					PhoneNumber: "2",
					Email:       "2",
					Raiting:     0,
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"name":"2","phone_number":"2","email":"2","raiting":0}`,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocks.NewMockUserRepo(ctrl)
			userService := service.NewUserService(userRepo)

			tt.mockBehavior(userRepo)

			service := service.Service{
				UserService: userService,
			}
			logger, _ := zap.NewProduction()
			handler := New(&service, nil, logger)

			r := gin.New()
			r.GET("/users/profile", handler.GetProfile)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/users/profile", nil)

			r.ServeHTTP(w, req)
			fmt.Println(w.Body.String())
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserRepo, user model.User)

	test := []struct {
		name         string
		body         string
		user         model.User
		mockBehavior mockBehavior
		expectedCode int
	}{
		{
			name: "update user",
			body: `{
				"phone_number": "+77777778",
				"email": "ripper@mail.ru"
			}`,
			user: model.User{
				PhoneNumber: "+77777778",
				Email:       "ripper@mail.ru",
			},
			mockBehavior: func(s *mocks.MockUserRepo, user model.User) {
				s.EXPECT().UpdateUserById(context.Background(), "", &user).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocks.NewMockUserRepo(ctrl)
			userService := service.NewUserService(userRepo)

			tt.mockBehavior(userRepo, tt.user)

			service := service.Service{
				UserService: userService,
			}
			logger, _ := zap.NewProduction()
			handler := New(&service, nil, logger)

			r := gin.New()
			r.PUT("/users/profile", handler.UpdateProfile)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/users/profile", bytes.NewBufferString(tt.body))

			r.ServeHTTP(w, req)
			fmt.Println(w.Body.String())
			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestDeleteProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserRepo)

	test := []struct {
		name         string
		mockBehavior mockBehavior
		expectedCode int
	}{
		{
			name: "delete user",
			mockBehavior: func(s *mocks.MockUserRepo) {
				s.EXPECT().DeleteUserById(context.Background(), "").Return(nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocks.NewMockUserRepo(ctrl)
			userService := service.NewUserService(userRepo)

			tt.mockBehavior(userRepo)

			service := service.Service{
				UserService: userService,
			}
			logger, _ := zap.NewProduction()
			handler := New(&service, nil, logger)

			r := gin.New()
			r.DELETE("/users", handler.DeleteUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/users", nil)

			r.ServeHTTP(w, req)
			fmt.Println(w.Body.String())
			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
