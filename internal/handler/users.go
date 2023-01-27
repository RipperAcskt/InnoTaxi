package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) Test(c *gin.Context) {
	users_id, _ := c.Get("users_id")
	if users_id != c.Param("id") {
		c.AbortWithStatus(403)
		return
	}
	c.Status(200)
}
