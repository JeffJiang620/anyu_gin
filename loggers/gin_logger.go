package loggers

import (
	"fmt"
	"github.com/JeffJiang620/anyu_logger/loggers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type LoggerConfig struct {
	SkipPaths []string
}

func Logger(notLogged ...string) gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
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
			loggers.Logger.Name("gin_logger").Info("",
				zap.String("ip", param.ClientIP),
				zap.String("proto", param.Request.Proto),
				zap.String("method", param.Method),
				zap.String("path", param.Path),
				zap.Int("status_code", param.StatusCode),
				zap.String("cost", fmt.Sprintf("%dms", param.Latency.Milliseconds())),
				zap.String("errMsg", param.ErrorMessage),
				zap.String("logType", "request"))
		}
	}
}
