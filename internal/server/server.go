package server

import (
	"context"

	"github.com/iryzzh/y-gophkeeper/internal/config"
	"github.com/iryzzh/y-gophkeeper/internal/server/web"
	"github.com/iryzzh/y-gophkeeper/internal/services/item"
	"github.com/iryzzh/y-gophkeeper/internal/services/token"
	"github.com/iryzzh/y-gophkeeper/internal/services/user"
)

type Server struct {
	webServerConfig *config.WebConfig
	debug           bool
	tokenSvc        *token.Service
	userSvc         *user.Service
	itemSvc         *item.Service
}

func NewServer(
	webServerConfig *config.WebConfig,
	tokenSvc *token.Service,
	userSvc *user.Service,
	itemSvc *item.Service,
	debug bool,
) *Server {
	return &Server{
		webServerConfig: webServerConfig,
		tokenSvc:        tokenSvc,
		userSvc:         userSvc,
		itemSvc:         itemSvc,
		debug:           debug,
	}
}

func (s *Server) Run(ctx context.Context) error {
	apiSrv := web.NewServer(
		s.webServerConfig.Network,
		s.webServerConfig.ServerAddress,
		s.webServerConfig.TLSCertPath,
		s.webServerConfig.TLSKeyPath,
		s.webServerConfig.EnableHTTPS,
		s.tokenSvc,
		s.userSvc,
		s.itemSvc,
		s.debug,
	)

	return apiSrv.Run(ctx)
}
