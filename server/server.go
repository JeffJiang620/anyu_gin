package server

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var errServerEngineNotInit = errors.New("server engine was not initialized, please call server.NewServer() to  initialize server")

type Server struct {
	engine *gin.Engine
	srv    *http.Server
}

func NewServer(mode string) *Server {
	gin.SetMode(mode)
	return &Server{
		engine: gin.New(),
	}
}

func (server *Server) Engine() *gin.Engine {
	return server.engine
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

func (server *Server) Start(addr ...string) error {

	if server.engine == nil {
		return errServerEngineNotInit
	}

	err := server.engine.SetTrustedProxies(nil)
	if err != nil {
		return err
	}

	address := resolveAddress(addr)

	server.srv = &http.Server{
		Addr:    address,
		Handler: server.engine,
	}

	err = server.srv.ListenAndServe()

	if err != nil {
		return err
	}
	return nil
}

func (server *Server) Stop(ctx context.Context) error {
	if server.srv == nil {
		return nil
	}
	return server.srv.Shutdown(ctx)
}

func (server *Server) WithOption(opts ...Option) (*Server, error) {
	if server.engine == nil {
		return nil, errServerEngineNotInit
	}
	for _, opt := range opts {
		err := opt.Apply(server)
		if err != nil {
			return nil, err
		}
	}

	return server, nil
}
