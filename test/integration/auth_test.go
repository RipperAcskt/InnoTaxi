package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/handler"
	"github.com/RipperAcskt/innotaxi/internal/repo/mongo"
	"github.com/RipperAcskt/innotaxi/internal/repo/postgres"
	"github.com/RipperAcskt/innotaxi/internal/repo/redis"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/go-playground/assert/v2"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func InitHandler() (*handler.Handler, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("config new failed: %w", err)
	}

	postgres, err := postgres.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres new failed: %w", err)
	}

	err = postgres.Migrate.Up()
	if err != migrate.ErrNoChange && err != nil {
		return nil, fmt.Errorf("migrate up failed: %w", err)
	}

	redis, err := redis.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("redis new failed: %w", err)
	}

	mongo, err := mongo.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("mongo new failed: %w", err)
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	writer := zapcore.AddSync(mongo)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	log := zap.New(core, zap.AddCaller())

	service := service.New(postgres, redis, cfg.SALT, cfg)
	return handler.New(service, cfg, log), nil
}

func TestSingUp(t *testing.T) {
	h, err := InitHandler()
	if err != nil {
		t.Errorf("init handler failed: %v", err)
	}

	test := []struct {
		name string
		body string
		code int
		err  error
	}{
		{
			name: "new user",
			body: `{"name": "Ivan", "phone_number": "+7455456", "email": "ripper@algsdh", "password": "12345"}`,
			code: http.StatusCreated,
			err:  nil,
		},
		{
			name: "existed user",
			body: `{"name": "Ivan", "phone_number": "+7455456", "email": "ripper@algsdh", "password": "12345"}`,
			code: http.StatusBadRequest,
			err:  service.ErrUserAlreadyExists,
		},
		{
			name: "empty body",
			body: `{}`,
			code: http.StatusBadRequest,
			err:  fmt.Errorf("EOF"),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			r := SetUpRouter()
			r.POST("/users/auth/sing-up", h.SingUp)

			req, _ := http.NewRequest("POST", "/users/auth/sing-up", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.IsEqual(tt.err, w.Body.String())
		})
	}
}

// func TestSingIn(t *testing.T) {
// 	h, err := InitHandler()
// 	if err != nil {
// 		t.Errorf("init handler failed: %v", err)
// 	}

// 	r := SetUpRouter()
// 	r.POST("/users/auth/sing-in", h.SingIn)
// 	t.Run("correct password", func(t *testing.T) {

// 		values := map[string]string{"phone_number": "+7455456", "password": "12345"}
// 		json_data, err := json.Marshal(values)

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		req, _ := http.NewRequest("POST", "/users/auth/sing-in", bytes.NewBuffer(json_data))
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusOK)
// 		}
// 	})

// 	t.Run("incorrect password", func(t *testing.T) {
// 		values := map[string]string{"phone_number": "+7455456", "password": "12345787979797979"}
// 		json_data, err := json.Marshal(values)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		req, _ := http.NewRequest("POST", "/users/auth/sing-in", bytes.NewBuffer(json_data))
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusForbidden {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusForbidden)
// 		}
// 	})

// }

// func TestRefresh(t *testing.T) {
// 	h, err := InitHandler()
// 	if err != nil {
// 		t.Errorf("init handler failed: %v", err)
// 	}

// 	r := SetUpRouter()
// 	r.GET("/users/auth/refresh", h.Refresh)
// 	r.POST("/users/auth/sing-in", h.SingIn)
// 	t.Run("without cookie", func(t *testing.T) {
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		req, _ := http.NewRequest("GET", "/users/auth/refresh", nil)
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusForbidden {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusForbidden)
// 		}
// 	})

// 	t.Run("with incorrect cookie", func(t *testing.T) {

// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		req, _ := http.NewRequest("GET", "/users/auth/refresh", nil)
// 		cookie := &http.Cookie{
// 			Name:   "refesh_token",
// 			Value:  "some_token",
// 			MaxAge: 300,
// 		}
// 		req.AddCookie(cookie)
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusForbidden {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusForbidden)
// 		}
// 	})

// 	t.Run("with incorrect signature cookie", func(t *testing.T) {
// 		values := map[string]string{"phone_number": "+7455456", "password": "12345"}
// 		json_data, err := json.Marshal(values)

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		req, _ := http.NewRequest("POST", "/users/auth/sing-in", bytes.NewBuffer(json_data))
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusOK)
// 		}
// 		cookies := w.Result().Cookies()

// 		req, _ = http.NewRequest("GET", "/users/auth/refresh", nil)
// 		w = httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		cookies[0].Value = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzgwMDgzMjMsInVzZXJfaWQiOjh9.YgzD0DlHj63RL1dw8l3IunpsxzY1b-JnIBPO35V9MY"

// 		req.AddCookie(cookies[0])
// 		w = httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusForbidden {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusForbidden)
// 		}
// 	})

// 	t.Run("with correct cookie", func(t *testing.T) {
// 		values := map[string]string{"phone_number": "+7455456", "password": "12345"}
// 		json_data, err := json.Marshal(values)

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		req, _ := http.NewRequest("POST", "/users/auth/sing-in", bytes.NewBuffer(json_data))
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusOK)
// 		}
// 		cookies := w.Result().Cookies()

// 		req, _ = http.NewRequest("GET", "/users/auth/refresh", nil)
// 		w = httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		req.AddCookie(cookies[0])
// 		w = httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusOK)
// 		}
// 	})

// }

// func TestLogout(t *testing.T) {
// 	h, err := InitHandler()
// 	if err != nil {
// 		t.Errorf("init handler failed: %v", err)
// 	}

// 	r := SetUpRouter()
// 	r.GET("/users/auth/logout", h.VerifyToken(), h.Logout)
// 	r.POST("/users/auth/sing-in", h.SingIn)
// 	t.Run("without access_token", func(t *testing.T) {
// 		req, _ := http.NewRequest("GET", "/users/auth/logout", nil)
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusUnauthorized {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusUnauthorized)
// 		}
// 	})

// 	t.Run("with access_token", func(t *testing.T) {
// 		values := map[string]string{"phone_number": "+7455456", "password": "12345"}
// 		json_data, err := json.Marshal(values)

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		req, _ := http.NewRequest("POST", "/users/auth/sing-in", bytes.NewBuffer(json_data))
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusOK)
// 		}
// 		token := make(map[string]string)
// 		json.Unmarshal(w.Body.Bytes(), &token)
// 		fmt.Println(token)

// 		req, _ = http.NewRequest("GET", "/users/auth/logout", nil)
// 		w = httptest.NewRecorder()
// 		req.Header.Add("Authorization", "Bearer "+token["access_token"])
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Errorf("got status %v, expected %v", w.Code, http.StatusOK)
// 		}
// 	})
// }
