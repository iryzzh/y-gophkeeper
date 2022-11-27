package web

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/iryzzh/gophkeeper/internal/services/item"

	"github.com/iryzzh/gophkeeper/internal/services/user"

	"github.com/iryzzh/gophkeeper/internal/services/token"

	"github.com/go-chi/chi/v5"
	v1 "github.com/iryzzh/gophkeeper/internal/server/web/api/v1"
	"github.com/iryzzh/gophkeeper/internal/tlsutil"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Server is an http server.
type Server struct {
	*chi.Mux
	network     string
	serverAddr  string
	tlsCertPath string
	tlsKeyPath  string
	enableHTTPS bool
	debug       bool
	tokenSvc    *token.Service
	userSvc     *user.Service
	itemSvc     *item.Service
}

// srvTimeout is the read and write timeout for the http server.
const srvTimeout = time.Second * 30

// NewServer returns a Server.
func NewServer(network, serverAddr, tlsCertPath, tlsKeyPath string, enableHTTPS bool, tokenSvc *token.Service,
	userSvc *user.Service, itemSvc *item.Service, debug bool) *Server {
	return &Server{
		network:     network,
		serverAddr:  serverAddr,
		enableHTTPS: enableHTTPS,
		tlsCertPath: tlsCertPath,
		tlsKeyPath:  tlsKeyPath,
		tokenSvc:    tokenSvc,
		userSvc:     userSvc,
		itemSvc:     itemSvc,
		debug:       debug,
	}
}

// Run configures and starts the http server.
func (s *Server) Run(ctx context.Context) error {
	listener, err := s.getListener()
	if err != nil {
		return err
	}
	defer func() {
		_ = listener.Close()
	}()

	s.Mux = chi.NewMux()
	s.registerMiddlewares()

	apiV1 := v1.NewAPI(s.tokenSvc, s.userSvc, s.itemSvc)
	apiV1.Register(s.Mux)

	srv := &http.Server{
		Handler:      s.Mux,
		ReadTimeout:  srvTimeout,
		WriteTimeout: srvTimeout,
	}

	srvType := "HTTP"
	if s.enableHTTPS {
		srvType = "HTTPS"
	}
	fmt.Printf("Starting the %s server on %s\n", srvType, listener.Addr())

	// serve in the background.
	serveError := make(chan error, 1)
	go func() {
		select {
		case serveError <- srv.Serve(listener):
		case <-ctx.Done():
		}
	}()

	// wait for stop or error signals.
	select {
	case <-ctx.Done():
		fmt.Printf("Shutting down %s server\n", srvType)
	case err = <-serveError:
		fmt.Printf("%s server error: %v\n", srvType, err)
		return err
	}

	timeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err = srv.Shutdown(timeout); errors.Is(err, timeout.Err()) { //nolint:contextcheck
		_ = srv.Close()
	}

	return nil
}

// getListener returns net.Listener. If enableHTTPS flag is
// set, it checks certificate validity. If the certificate is
// invalid, a new one is generated.
func (s *Server) getListener() (net.Listener, error) {
	if s.enableHTTPS {
		cert, err := tls.LoadX509KeyPair(s.tlsCertPath, s.tlsKeyPath)
		if err != nil {
			fmt.Println("Loading certificate: ", err)
			fmt.Println("Creating a new certificate")

			cert, err = tlsutil.NewCertificate(s.tlsCertPath, s.tlsKeyPath, "gophkeeper")
			if err != nil {
				return nil, err
			}
		}

		tlsConfig := &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
		}

		return tls.Listen(s.network, s.serverAddr, tlsConfig)
	}

	return net.Listen(s.network, s.serverAddr)
}
