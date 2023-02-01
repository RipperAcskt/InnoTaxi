package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/RipperAcskt/innotaxi/config"
	_ "github.com/RipperAcskt/innotaxi/docs"
	"github.com/RipperAcskt/innotaxi/internal/service"
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

	fmt.Println(swaggerFiles.Handler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	users := router.Group("/users")

	auth := users.Group("/auth")
	auth.POST("sing-up", h.singUp)
	auth.POST("sing-in", h.singIn)
	auth.POST("refresh", h.Refresh)

	users.GET("/profile/:id", VerifyToken(h.cfg), h.GetProfile)
	users.PUT("/profile/:id", VerifyToken(h.cfg), h.UpdateProfile)
	users.DELETE("/:id", VerifyToken(h.cfg), h.DeleteUser)
	auth.GET("refresh", h.Refresh)
	auth.GET("logout", h.VerifyToken(), h.Logout)
	return router
}
