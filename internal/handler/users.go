package handler

import (
	"net/http"

	"github.com/RipperAcskt/innotaxi/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetProfile(c *gin.Context) {
	user, err := h.s.GetProfile(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":         user.Name,
		"phone_number": user.PhoneNumber,
		"email":        user.Email,
		"raiting":      user.Raiting,
	})
}

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
