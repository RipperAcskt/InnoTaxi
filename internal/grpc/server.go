package grpc_server

import (
	"context"
	"fmt"
	"net"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"google.golang.org/grpc"
)

type Server struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{cfg}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.cfg.GRPC_HOST)

	if err != nil {
		return fmt.Errorf("listen failed: %w", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	service.RegisterAuthServiceServer(grpcServer, s)
	grpcServer.Serve(listener)

	return nil
}

func (s *Server) GetJWT(c context.Context, params *service.Params) (*service.Response, error) {
	tokenParams := service.TokenParams{
		DriverID: params.DriverID,
		Type:     params.Type,
	}

	token, err := service.NewToken(tokenParams, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("new token failed: %w", err)
	}

	response := &service.Response{
		AccessToken:  token.Access,
		RefreshToken: token.RT,
	}
	return response, nil
}
