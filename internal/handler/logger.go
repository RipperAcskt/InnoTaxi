package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *Handler) Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uuid := uuid.New()
		c.Set("start", start)
		c.Set("uuid", uuid)
		c.Next()

		h.log.Info("request", zap.String("url", c.Request.URL.Path), zap.String("method", c.Request.Method), zap.Any("uuid", uuid), zap.String("time", time.Since(start).String()))
	}
}

func getTimeAndId(c *gin.Context) (start time.Time, uuid any) {
	t, _ := c.Get("start")
	start, _ = t.(time.Time)
	uuid, _ = c.Get("uuid")
	return
}
