package handler

import (
	"fmt"
	"net/http"

	"github.com/RipperAcskt/innotaxi/internal/model"
	"github.com/gin-gonic/gin"
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
	user, err := h.s.GetProfile(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get profile failed: %w", err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
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
	var user model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.s.UpdateProfile(c.Request.Context(), c.Param("id"), &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(200)
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
	err := h.s.DeleteUser(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(200)
}
