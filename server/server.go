package server

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var errServerEngineNotInit = errors.New("server engine was not initialized, please call server.NewServer() to  initialize server")

type Server struct {
	engine *gin.Engine
}

func NewServer() *Server {
	return &Server{
		engine: gin.New(),
	}
}

func (server *Server) Engine() *gin.Engine {
	return server.engine
}

func (server *Server) Start(addr ...string) error {

	if server.engine == nil {
		return errServerEngineNotInit
	}

	err := server.engine.SetTrustedProxies(nil)
	if err != nil {
		return err
	}
	err = server.engine.Run(addr...)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) Stop() error {
	return nil
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
