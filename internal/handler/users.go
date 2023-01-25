package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) Test(c *gin.Context) {
	c.Status(200)
}
