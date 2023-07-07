package server

import (
	"errors"
	"strings"
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
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	MaxAge           int
}

func (opt Cors) Apply(server *Server) error {

	allowOrigins := strings.Split(opt.AllowOrigins, ",")
	allowMethods := strings.Split(opt.AllowMethods, ",")
	allowHeaders := strings.Split(opt.AllowHeaders, ",")
	maxAge := time.Duration(opt.MaxAge) * time.Minute

	c := cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: opt.AllowCredentials,
		MaxAge:           maxAge,
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
	Locale     string
	Validators []validators.Validator
}

func (opt Validators) Apply(server *Server) error {
	if !checkLocaleSupported(opt.Locale) {
		return errNotSupportedLocale
	}

	return validators.RegisterValidator(opt.Locale, opt.Validators...)
}

type Routers struct {
	Base    string
	Routers []routers.Router
}

func (opt Routers) Apply(server *Server) error {
	routers.CombineRouters(server.engine, opt.Base, opt.Routers...)
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

type Middlewares []gin.HandlerFunc

func (m Middlewares) Apply(server *Server) error {
	for _, handler := range m {
		server.Engine().Use(handler)
	}

	return nil
}
