package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/RipperAcskt/innotaxi/internal/model"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// @Summary get user profile
// @Tags user
// @Param id path int true "user's id"
// @Produce json
// @Success 200 {object} model.User
// @Failure 401 {object} error "error: err"
// @Failure 403 {object} error "error: err"
// @Failure 500 {object} error "error: err"
// @Router /users/profile/{id} [GET]
// @Security Bearer
func (h *Handler) GetProfile(c *gin.Context) {
	start := time.Now()
	uuid := uuid.New()

	user, err := h.s.GetProfile(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) {
			h.log.Info("/users/profile/{id}", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.String("time", time.Since(start).String()))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.log.Error("/users/profile/{id}", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.Error(fmt.Errorf("get profile failed: %w", err)), zap.String("time", time.Since(start).String()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get profile failed: %w", err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
	h.log.Info("/users/profile/{id}", zap.String("method", "GET"), zap.Any("uuid", uuid), zap.String("time", time.Since(start).String()))
}

// @Summary update user profile
// @Tags user
// @Param input body model.User false "rows to update"
// @Param id path int true "user's id"
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Failure 401 {object} error "error: err"
// @Failure 403 {object} error "error: err"
// @Failure 500 {object} error "error: err"
// @Router /users/profile/{id} [PUT]
// @Security Bearer
func (h *Handler) UpdateProfile(c *gin.Context) {
	start := time.Now()
	uuid := uuid.New()
	var user model.User

	if err := c.BindJSON(&user); err != nil {
		h.log.Info("/users/profile/{id}", zap.String("method", "PUT"), zap.Any("uuid", uuid), zap.String("time", time.Since(start).String()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.s.UpdateProfile(c.Request.Context(), c.Param("id"), &user)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) {
			h.log.Info("/users/profile/{id}", zap.String("method", "PUT"), zap.Any("uuid", uuid), zap.String("time", time.Since(start).String()))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		h.log.Error("/users/profile/{id}", zap.String("method", "PUT"), zap.Any("uuid", uuid), zap.Error(fmt.Errorf("update profile failed: %w", err)), zap.String("time", time.Since(start).String()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(200)
	h.log.Info("/users/profile/{id}", zap.String("method", "PUT"), zap.Any("uuid", uuid), zap.String("time", time.Since(start).String()))
}

// @Summary delete user
// @Tags user
// @Param id path int false "user's id to delete"
// @Accept json
// @Success 200
// @Failure 401 {object} error "error: err"
// @Failure 403 {object} error "error: err"
// @Failure 500 {object} error "error: err"
// @Router /users/{id} [DELETE]
// @Security Bearer
func (h *Handler) DeleteUser(c *gin.Context) {
	start := time.Now()
	uuid := uuid.New()

	err := h.s.DeleteUser(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) {
			h.log.Info("/users/{id}", zap.String("method", "DELETE"), zap.Any("uuid", uuid), zap.String("time", time.Since(start).String()))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.log.Error("/users/{id}", zap.String("method", "DELETE"), zap.Any("uuid", uuid), zap.Error(fmt.Errorf("delete user failed: %w", err)), zap.String("time", time.Since(start).String()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(200)
	h.log.Info("/users/{id}", zap.String("method", "DELETE"), zap.Any("uuid", uuid), zap.String("time", time.Since(start).String()))
}
