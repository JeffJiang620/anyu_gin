package loggers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

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
