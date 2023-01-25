package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/RipperAcskt/innotaxi/internal/service"
)

func (h *Handler) singUp(c *gin.Context) {
	var user service.UserSingUp

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx := context.Background()
	err := h.s.CreateUser(ctx, user)
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

	ctx := context.Background()
	token, err := h.s.SingIn(ctx, user)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) || errors.Is(err, service.ErrIncorrectPassword) {
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

	exp, err := strconv.Atoi(os.Getenv("RTEXP"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("atoi rt failed: %w", err),
		})
	}
	jwtExp := int((time.Duration(exp) * time.Hour * 24).Seconds())
	c.SetCookie("refresh_token", token.RT, jwtExp, "/users/auth", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": token.Access,
	})
}

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := make(map[string]string)

		if err := c.BindJSON(&accessToken); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ok, err := service.Verify(accessToken["access_token"])
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Errorf("wrong signature").Error(),
			})
			return
		}
		c.Next()

	}
}

func (h *Handler) Refresh(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": fmt.Errorf("get refresh token failed").Error(),
		})
	}

	ok, err := service.Verify(rt)
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("wrong signature").Error(),
		})
		return
	}

	token, err := service.NewToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("wrong signature").Error(),
		})
		return
	}

	exp, err := strconv.Atoi(os.Getenv("RTEXP"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("atoi rt failed: %w", err),
		})
	}
	jwtExp := int((time.Duration(exp) * time.Hour * 24).Seconds())
	c.SetCookie("refresh_token", token.RT, jwtExp, "/users/auth", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": token.Access,
	})
}
