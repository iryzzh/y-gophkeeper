package server

import (
	"context"

	"github.com/iryzzh/gophkeeper/internal/services/item"

	"github.com/iryzzh/gophkeeper/internal/services/user"

	"github.com/iryzzh/gophkeeper/internal/services/token"

	"github.com/iryzzh/gophkeeper/internal/config"
	"github.com/iryzzh/gophkeeper/internal/server/web"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	webServerConfig *config.WebServerConfig
	debug           bool
	tokenSvc        *token.Service
	userSvc         *user.Service
	itemSvc         *item.Service
}

func NewServer(
	webServerConfig *config.WebServerConfig,
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
	g, childCtx := errgroup.WithContext(ctx)

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

	g.Go(func() error {
		return apiSrv.Run(childCtx)
	})

	return g.Wait()
}
