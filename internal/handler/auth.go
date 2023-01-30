package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/service"
)

// SingUp gosoc
// @Summary registarte user
// @Tags auth
// @Accept json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /users/auth/sing-up [POST]
func (h *Handler) singUp(c *gin.Context) {
	var user service.UserSingUp

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.s.SingUp(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (h *Handler) singIn(c *gin.Context) {
	var user service.UserSingIn

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.s.SingIn(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) || errors.Is(err, service.ErrIncorrectPassword) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

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
}

func VerifyToken(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := strings.Split(c.GetHeader("Authorization"), " ")[1]

		ok, id, err := service.Verify(accessToken, cfg)
		if err != nil {
			if errors.Is(err, service.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}
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
		c.Set("users_id", fmt.Sprintf("%v", id))
		c.Next()

	}
}

func (h *Handler) Refresh(c *gin.Context) {
	refresh, err := c.Cookie("refresh_token")

	if err != nil {
		if err == http.ErrNoCookie {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Errorf("bad refresh token").Error(),
			})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	ok, id, err := service.Verify(refresh, h.cfg)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("verify rt failed: %w", err).Error(),
		})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": fmt.Errorf("wrong signature").Error(),
		})
		return
	}

	token, err := service.NewToken(id, h.cfg)
	if err != nil {
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
}
