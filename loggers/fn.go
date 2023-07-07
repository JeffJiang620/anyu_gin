package loggers

import (
	"fmt"

	"github.com/anyufly/logger/loggers"
	"github.com/anyufly/stack_err/stackerr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogRequestErr(context *gin.Context, err error) {
	path := context.Request.URL.Path
	ip := context.ClientIP()
	method := context.Request.Method

	if e, ok := err.(stackerr.ErrorWithStack); ok {
		loggers.Logger.Name("request_error").Error("",
			zap.String("ip", ip),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("errorSource", fmt.Sprintf("%s:%d", e.File(), e.Line())),
			zap.String("errMsg", e.Error()))
	} else {
		loggers.Logger.Name("request_error").Error("",
			zap.String("ip", ip),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("errMsg", err.Error()))
	}
}
