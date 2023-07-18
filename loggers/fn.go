package loggers

import (
	"fmt"

	"github.com/anyufly/stack_err/stackerr"
	"github.com/gin-gonic/gin"
)

func LogRequestErr(context *gin.Context, err error) {
	logger := Default(context)

	if logger != nil {
		path := context.Request.URL.Path
		ip := context.ClientIP()
		method := context.Request.Method

		if e, ok := err.(stackerr.ErrorWithStack); ok {
			logger.Name("request_error").Error("",
				"ip", ip,
				"method", method,
				"path", path,
				"errorSource", fmt.Sprintf("%s:%d", e.File(), e.Line()),
				"errMsg", e.Error())
		} else {
			logger.Name("request_error").Error("",
				"ip", ip,
				"method", method,
				"path", path,
				"errMsg", err.Error())
		}
	}

}
