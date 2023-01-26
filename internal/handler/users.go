package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) Test(c *gin.Context) {
	id, _ := c.Get("id")
	if id != c.Param("id") {
		c.AbortWithStatus(403)
		return
	}
	c.Status(200)
}
