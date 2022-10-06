package server

import "github.com/gin-gonic/gin"

type Server struct {
	mode   string
	engine *gin.Engine
}

func (server *Server) Engine() *gin.Engine {
	return server.engine
}

func (server *Server) Mode() string {
	return server.mode
}

func (server *Server) Start(addr ...string) error {

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
	for _, opt := range opts {
		err := opt.Apply(server)
		if err != nil {
			return nil, err
		}
	}

	return server, nil
}
