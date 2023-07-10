package loggers

import (
	"github.com/anyufly/logger/loggers"
	"github.com/gin-gonic/gin"
)

type GinLogger interface {
	Name(name string) GinLogger
	Info(msg string, keyAndValues ...interface{})
	Debug(msg string, keyAndValues ...interface{})
	Warn(msg string, keyAndValues ...interface{})
	Error(msg string, keyAndValues ...interface{})
}

const loggerKey = "_gin_logger"

func Default(ctx *gin.Context) GinLogger {
	logger, ok := ctx.Get(loggerKey)

	if !ok {
		return nil
	}

	return logger.(GinLogger)
}

type ginLogger struct {
	logger *loggers.CommonLogger
}

func New(logger *loggers.CommonLogger) GinLogger {
	return &ginLogger{
		logger: logger,
	}
}

func (g *ginLogger) Name(name string) GinLogger {
	newLogger := g.logger.Name(name)
	return New(newLogger)

}

func (g *ginLogger) Info(msg string, keyAndValues ...interface{}) {
	g.logger.Sugar().Infow(msg, keyAndValues...)
}

func (g *ginLogger) Debug(msg string, keyAndValues ...interface{}) {
	g.logger.Sugar().Debugw(msg, keyAndValues...)
}

func (g *ginLogger) Warn(msg string, keyAndValues ...interface{}) {
	g.logger.Sugar().Warnw(msg, keyAndValues...)
}

func (g *ginLogger) Error(msg string, keyAndValues ...interface{}) {
	g.logger.Sugar().Errorw(msg, keyAndValues...)
}
