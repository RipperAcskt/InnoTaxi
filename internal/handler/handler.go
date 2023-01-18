package handler

import (
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{s}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	auth := router.Group("/users")
	auth.POST("singUp", h.singUp)

	return router
}
