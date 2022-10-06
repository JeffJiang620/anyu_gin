package server

import (
	"errors"
	"time"

	"github.com/anyufly/gin_common/routers"
	"github.com/anyufly/gin_common/validators"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Option interface {
	Apply(*Server) error
}

var supportedMode = []string{
	gin.DebugMode,
	gin.ReleaseMode,
	gin.TestMode,
}

var errNotSupportedMode = errors.New("mode not supported")

func checkModeSupported(mode string) bool {
	for _, sm := range supportedMode {
		if mode == sm {
			return true
		}
	}

	return false
}

type Mode string

func (opt Mode) Apply(server *Server) error {
	sMode := string(opt)
	if !checkModeSupported(sMode) {
		return errNotSupportedMode
	}

	server.mode = sMode
	gin.SetMode(sMode)
	return nil
}

type Logger struct {
	LoggerHandler gin.HandlerFunc
}

func (opt Logger) Apply(server *Server) error {
	server.engine.Use(opt.LoggerHandler)
	return nil
}

type Recover struct {
	RecoverHandler gin.HandlerFunc
}

func (opt Recover) Apply(server *Server) error {
	server.engine.Use(opt.RecoverHandler)
	return nil
}

type Cors struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           time.Duration
}

func (opt Cors) Apply(server *Server) error {
	c := cors.Config{
		AllowOrigins:     opt.AllowOrigins,
		AllowMethods:     opt.AllowMethods,
		AllowHeaders:     opt.AllowHeaders,
		AllowCredentials: opt.AllowCredentials,
		MaxAge:           opt.MaxAge,
	}
	server.engine.Use(cors.New(c))
	return nil
}

type PProfEnable bool

func (opt PProfEnable) Apply(server *Server) error {
	if bool(opt) && (gin.IsDebugging() || gin.Mode() == gin.TestMode) {
		pprof.Register(server.engine)
	}

	return nil
}

var errNotSupportedLocale = errors.New("locale not supported")

var supportedLocale = []string{
	"en", "zh",
}

func checkLocaleSupported(locale string) bool {
	for _, sl := range supportedLocale {
		if locale == sl {
			return true
		}
	}

	return false
}

type Validators struct {
	locale     string
	validators []validators.Validator
}

func (opt Validators) Apply(server *Server) error {
	if !checkLocaleSupported(opt.locale) {
		return errNotSupportedLocale
	}

	return validators.RegisterValidator(opt.locale, opt.validators...)
}

type Routers struct {
	base    string
	routers []routers.Router
}

func (opt Routers) Apply(server *Server) error {
	routers.CombineRouters(server.engine, opt.base, opt.routers...)
	return nil
}

type SwaggerEnable bool

func (opt SwaggerEnable) Apply(server *Server) error {
	if bool(opt) && (gin.IsDebugging() || gin.Mode() == gin.TestMode) {
		server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, func(c *ginSwagger.Config) {
			c.DefaultModelsExpandDepth = -1
		}))
	}
	return nil
}
