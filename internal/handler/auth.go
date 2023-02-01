package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

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
	var user service.UserSingUp

	if err := c.BindJSON(&user); err != nil {
		h.log.WithFields(logrus.Fields{
			"package":  "handler",
			"method":   "singUp",
			"function": "BindJSON",
			"error":    err,
		}).Warning("bind json failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.s.SingUp(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			h.log.WithFields(logrus.Fields{
				"package":  "handler",
				"method":   "singUp",
				"function": "h.s.SingUp",
				"error":    err,
				"user":     user,
			}).Warning("service sing up failed")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.log.WithFields(logrus.Fields{
			"package":  "handler",
			"method":   "singUp",
			"function": "h.s.SingUp",
			"error":    err,
			"user":     user,
		}).Error("service sing up failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
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
	var user service.UserSingIn

	if err := c.BindJSON(&user); err != nil {
		h.log.WithFields(logrus.Fields{
			"package":  "handler",
			"method":   "singIn",
			"function": "BindJSON",
			"error":    err,
		}).Warning("bind json failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.s.SingIn(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) || errors.Is(err, service.ErrIncorrectPassword) {
			h.log.WithFields(logrus.Fields{
				"package":  "handler",
				"method":   "singIn",
				"function": "h.s.SingIp",
				"error":    err,
				"user":     user,
			}).Warning("service sing in failed")
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.log.WithFields(logrus.Fields{
			"package":  "handler",
			"method":   "singIn",
			"function": "h.s.SingIn",
			"error":    err,
			"user":     user,
		}).Error("service sing in failed")
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

func (h *Handler) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.Split(c.GetHeader("Authorization"), " ")
		if len(token) < 2 {
			h.log.WithFields(logrus.Fields{
				"package":  "handler",
				"method":   "VerifyToken",
				"function": "len",
				"error":    "can't get access token",
				"token":    token,
			}).Warning("access token required")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Errorf("access token required").Error(),
			})
			return
		}
		accessToken := token[1]

		ok, id, err := service.Verify(accessToken, h.cfg)
		if err != nil {
			if errors.Is(err, service.ErrTokenExpired) {
				h.log.WithFields(logrus.Fields{
					"package":  "handler",
					"method":   "VerifyToken",
					"function": "service.Verify",
					"error":    err,
					"token":    accessToken,
				}).Warning("service sing in failed")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}

			h.log.WithFields(logrus.Fields{
				"package":  "handler",
				"method":   "VerifyToken",
				"function": "service.Verify",
				"error":    err,
				"token":    accessToken,
			}).Error("service verify failed")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Errorf("verify failed: %w", err).Error(),
			})
			return
		}
		if !ok {
			h.log.WithFields(logrus.Fields{
				"package":  "handler",
				"method":   "VerifyToken",
				"function": "service.Verify",
				"error":    "bad access token",
				"token":    accessToken,
			}).Warning("service verify failed")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Errorf("wrong signature").Error(),
			})
			return
		}

		if !h.s.CheckToken(accessToken) {
			h.log.WithFields(logrus.Fields{
				"package":  "handler",
				"method":   "VerifyToken",
				"function": "h.s.CheckToken",
				"error":    "user logouted",
				"token":    accessToken,
			}).Warning("service check token failed")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if fmt.Sprint(id) != c.Param("id") {
			h.log.WithFields(logrus.Fields{
				"package":  "handler",
				"method":   "VerifyToken",
				"error":    "identifiers do not equel",
				"id":       id,
				"param id": id,
			}).Warning("compare identifiers failed")
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
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			h.log.WithFields(logrus.Fields{
				"package":  "handler",
				"method":   "refresh",
				"function": "c.Cookie",
				"error":    err,
			}).Warning("get cookie failed")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Errorf("bad refresh token").Error(),
			})
			return
		} else {
			h.log.WithFields(logrus.Fields{
				"package":  "handler",
				"method":   "refresh",
				"function": "c.Cookie",
				"error":    err,
			}).Error("get cookie failed")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	ok, id, err := service.Verify(refresh, h.cfg)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			h.log.WithFields(logrus.Fields{
				"package":  "handler",
				"method":   "refresh",
				"function": "service.Verify",
				"error":    err,
				"token":    refresh,
			}).Warning("verify failed")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.log.WithFields(logrus.Fields{
			"package":  "handler",
			"method":   "refresh",
			"function": "service.Verify",
			"error":    err,
			"token":    refresh,
		}).Error("verify failed")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("verify rt failed: %w", err).Error(),
		})
		return
	}
	if !ok {
		h.log.WithFields(logrus.Fields{
			"package":  "handler",
			"method":   "refresh",
			"function": "service.Verify",
			"error":    err,
			"token":    refresh,
		}).Warning("wrong signarure")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": fmt.Errorf("wrong signature").Error(),
		})
		return
	}

	token, err := service.NewToken(id, h.cfg)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"package":  "handler",
			"method":   "refresh",
			"function": "service.NewToken",
			"error":    err,
			"id":       id,
		}).Warning("new token failed")
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

// @Summary refresh access token
// @Tags auth
// @Produce json
// @Success 200
// @Failure 401 {object} error  "error: err"
// @Failure 403 {object} error  "error: err"
// @Failure 500 {object} error  "error: err"
// @Router /users/auth/logout [GET]
// @Security Bearer
func (h *Handler) Logout(c *gin.Context) {
	id, ok := c.Get("id")
	if !ok {
		h.log.WithFields(logrus.Fields{
			"package":  "handler",
			"method":   "Logout",
			"function": "c.Get",
			"error":    "can't get id",
			"id":       id,
		}).Warning("get id failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "can't get id",
		})
		return
	}

	exp := time.Duration(h.cfg.ACCESS_TOKEN_EXP) * time.Minute
	token := strings.Split(c.GetHeader("Authorization"), " ")
	if len(token) < 2 {
		h.log.WithFields(logrus.Fields{
			"package":  "handler",
			"method":   "Logout",
			"function": "len",
			"error":    "can't get access token",
			"token":    token,
		}).Warning("access token required")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Errorf("access token required").Error(),
		})
		return
	}
	accessToken := token[1]

	err := h.s.Logout(id.(string), accessToken, exp)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"package":  "handler",
			"method":   "Logout",
			"function": "h.s.Logout",
			"error":    err,
			"token":    accessToken,
		}).Warning("get id failed")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.SetCookie("refresh_token", "", time.Now().Second(), "/users/auth", "", false, true)
	c.Status(http.StatusOK)
}
