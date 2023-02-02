package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/RipperAcskt/innotaxi/internal/service"
)

// @Summary registrate user
// @Tags auth
// @Param user body service.UserSingUp true "account info"
// @Accept json
// @Success 200
// @Failure 400 {object} error "error: err"
// @Failure 500 {object} error "error: err"
// @Router /users/auth/sing-up [POST]
func (h *Handler) singUp(c *gin.Context) {
	start := time.Now()
	uuid := uuid.New()
	var user service.UserSingUp

	if err := c.BindJSON(&user); err != nil {
		h.log.Info("/users/auth/sing-up", zap.String("method", "POST"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.s.SingUp(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			h.log.Info("/users/auth/sing-up", zap.String("method", "POST"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.log.Error("/users/auth/sing-up", zap.String("method", "POST"), zap.Any("uuid", uuid), zap.Error(fmt.Errorf("service sing up failed: %w", err)), zap.Duration("time", time.Since(start)))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	h.log.Info("/users/auth/sing-up", zap.String("method", "POST"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
}

// @Summary user authentication
// @Tags auth
// @Param input body service.UserSingIn true "phone number and password"
// @Accept json
// @Produce json
// @Success 200 {object} string "access_token: token"
// @Failure 403 {object} error "error: err"
// @Failure 500 {object} error "error: err"
// @Router /users/auth/sing-in [POST]
func (h *Handler) singIn(c *gin.Context) {
	start := time.Now()
	uuid := uuid.New()
	var user service.UserSingIn

	if err := c.BindJSON(&user); err != nil {
		h.log.Info("/users/auth/sing-in", zap.String("method", "POST"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.s.SingIn(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) || errors.Is(err, service.ErrIncorrectPassword) {
			h.log.Info("/users/auth/sing-in", zap.String("method", "POST"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.log.Error("/users/auth/sing-in", zap.String("method", "POST"), zap.Any("uuid", uuid), zap.Error(fmt.Errorf("service sing in failed: %w", err)), zap.Duration("time", time.Since(start)))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	exp := int((time.Duration(h.cfg.REFRESH_TOKEN_EXP) * time.Hour * 24).Seconds())
	c.SetCookie("refresh_token", token.RT, exp, "/users/auth", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"access_token": token.Access,
	})
	h.log.Info("/users/auth/sing-in", zap.String("method", "POST"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
}

func (h *Handler) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.Split(c.GetHeader("Authorization"), " ")
		if len(token) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Errorf("access token required").Error(),
			})
			return
		}
		accessToken := token[1]

		ok, id, err := service.Verify(accessToken, h.cfg)
		if err != nil {
			if errors.Is(err, service.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}

			h.log.Error("service verify failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Errorf("verify failed: %w", err).Error(),
			})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Errorf("wrong signature").Error(),
			})
			return
		}

		if !h.s.CheckToken(accessToken) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if fmt.Sprint(id) != c.Param("id") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()

	}
}

// @Summary refresh access token
// @Tags auth
// @Produce json
// @Success 200 {object} string "accept token: token"
// @Failure 401 {object} error  "error: err"
// @Failure 403 {object} error  "error: err"
// @Failure 500 {object} error  "error: err"
// @Router /users/auth/refresh [GET]
func (h *Handler) Refresh(c *gin.Context) {
	start := time.Now()
	uuid := uuid.New()

	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			h.log.Info("/users/auth/refresh", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Errorf("bad refresh token").Error(),
			})
			return
		}
		h.log.Error("/users/auth/refresh", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Error(fmt.Errorf("get cookie failed: %w", err)), zap.Duration("time", time.Since(start)))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ok, id, err := service.Verify(refresh, h.cfg)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			h.log.Info("/users/auth/refresh", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.log.Error("/users/auth/refresh", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Error(fmt.Errorf("verify failed: %w", err)), zap.Duration("time", time.Since(start)))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("verify rt failed: %w", err).Error(),
		})
		return
	}
	if !ok {
		h.log.Info("/users/auth/refresh", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": fmt.Errorf("wrong signature").Error(),
		})
		return
	}

	token, err := service.NewToken(id, h.cfg)
	if err != nil {
		h.log.Info("/users/auth/refresh", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": fmt.Errorf("wrong signature").Error(),
		})
		return
	}

	exp := int((time.Duration(h.cfg.REFRESH_TOKEN_EXP) * time.Hour * 24).Seconds())
	c.SetCookie("refresh_token", token.RT, exp, "/users/auth", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"access_token": token.Access,
	})
	h.log.Info("/users/auth/refresh", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
}

// @Summary logout user
// @Tags auth
// @Accept json
// @Success 200
// @Failure 401 {object} error  "error: err"
// @Failure 403 {object} error  "error: err"
// @Failure 500 {object} error  "error: err"
// @Router /users/auth/logout [GET]
// @Security Bearer
func (h *Handler) Logout(c *gin.Context) {
	start := time.Now()
	uuid := uuid.New()
	id, ok := c.Get("id")
	if !ok {
		h.log.Info("/users/auth/logout", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "can't get id",
		})
		return
	}

	exp := time.Duration(h.cfg.ACCESS_TOKEN_EXP) * time.Minute
	token := strings.Split(c.GetHeader("Authorization"), " ")
	if len(token) < 2 {
		h.log.Info("/users/auth/logout", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Errorf("access token required").Error(),
		})
		return
	}
	accessToken := token[1]

	err := h.s.Logout(id.(string), accessToken, exp)
	if err != nil {
		h.log.Error("/users/auth/logout", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Error(fmt.Errorf("logout failed: %w", err)), zap.Duration("time", time.Since(start)))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.SetCookie("refresh_token", "", time.Now().Second(), "/users/auth", "", false, true)
	c.Status(http.StatusOK)
	h.log.Info("/users/auth/logout", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Duration("time", time.Since(start)))
}
