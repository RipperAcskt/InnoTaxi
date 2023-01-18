package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RipperAcskt/innotaxi/internal/model"
)

func (h *Handler) singUp(c *gin.Context) {
	var user model.UserSingUp

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "bad request",
			"error":  fmt.Sprint(err),
		})
		return
	}

	err := h.s.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal server error",
			"error":  fmt.Sprint(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
