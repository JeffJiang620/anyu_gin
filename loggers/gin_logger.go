package loggers

import (
	"fmt"
	"time"

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

type LoggerConfig struct {
	Logger    GinLogger
	SkipPaths []string
}

func Logger(logger GinLogger, notLogged ...string) gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		Logger:    logger,
		SkipPaths: notLogged,
	})
}

func LoggerWithConfig(conf LoggerConfig) gin.HandlerFunc {
	notLogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		c.Set(loggerKey, conf.Logger)
		// Start time
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}
			// Stop time
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			conf.Logger.Info("",
				"ip", param.ClientIP,
				"proto", param.Request.Proto,
				"method", param.Method,
				"path", param.Path,
				"status_code", param.StatusCode,
				"cost", fmt.Sprintf("%d ms", param.Latency.Milliseconds()),
				"errMsg", param.ErrorMessage,
				"logType", "request")
		}
	}
}

func Default(ctx *gin.Context) GinLogger {
	logger, ok := ctx.Get(loggerKey)

	if !ok {
		return nil
	}

	return logger.(GinLogger)
}
