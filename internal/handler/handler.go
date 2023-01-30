package handler

import (
	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	s   *service.Service
	cfg *config.Config
}

func New(s *service.Service, cfg *config.Config) *Handler {
	return &Handler{s, cfg}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	users := router.Group("/users")

	auth := users.Group("/auth")
	auth.POST("sing-up", h.singUp)
	auth.POST("sing-in", h.singIn)
	auth.GET("refresh", h.Refresh)
	auth.GET("logout", h.VerifyToken(), h.Logout)
	return router
}
