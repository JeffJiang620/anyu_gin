package middlewares

import (
	"errors"
	ginRender "github.com/gin-gonic/gin/render"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/anyufly/gin_common/loggers"
	"github.com/anyufly/gin_common/renders"
	"github.com/anyufly/gin_common/response"
	"github.com/gin-gonic/gin"
)

const (
	phaseBefore = 1
	phaseAfter  = 2
)

var ErrUnsupportedMiddlewareReturnType = errors.New("unsupported middleware return type")

func processMiddlewareFunc(phase int, middleware IMiddleWare, ctx *gin.Context) (goOn bool) {
	var allow bool
	var processFunc func(context *gin.Context) interface{}
	switch phase {
	case phaseBefore:
		processFunc = middleware.Before
		allow = !middleware.DeniedBeforeAbortContext()
	case phaseAfter:
		processFunc = middleware.After
		allow = middleware.AllowAfterAbortContext()
	}

	data := processFunc(ctx)

	if data == nil {
		goOn = true
		return
	}

	switch r := data.(type) {
	case renders.ErrorRender:
		err := r.Err()
		loggers.LogRequestErr(ctx, err)
		if allow {
			r.Render(ctx)
			ctx.Abort()
			return
		}
	case renders.Render:
		if allow {
			r.Render(ctx)
			ctx.Abort()
			return
		}

	case ginRender.Render:
		if allow {
			ctx.Render(http.StatusOK, r)
			ctx.Abort()
			return
		}
	case validator.ValidationErrors:
		er := response.ParameterError.WithErr(r)
		er.Render(ctx)
		ctx.Abort()
		return
	case error:
		loggers.LogRequestErr(ctx, r)
		if allow {
			er := response.UnknownError.WithErr(r)
			er.Render(ctx)
			ctx.Abort()
			return
		}
	default:
		panic(ErrUnsupportedMiddlewareReturnType)
	}

	goOn = true
	return
}

func MiddlewareHandler(middleware IMiddleWare) gin.HandlerFunc {
	return func(context *gin.Context) {
		if ok := processMiddlewareFunc(phaseBefore, middleware, context); !ok {
			return
		}
		context.Next()

		processMiddlewareFunc(phaseAfter, middleware, context)
	}
}
