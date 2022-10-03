package controllers

import (
	"github.com/JeffJiang620/anyu_gin/loggers"
	"github.com/JeffJiang620/anyu_gin/renders"
	"github.com/JeffJiang620/anyu_gin/response"
	"github.com/gin-gonic/gin"
	ginRender "github.com/gin-gonic/gin/render"
	"net/http"
)

type ContextOption interface {
	Apply(ctx *gin.Context)
}

func WithOption(ctx *gin.Context, opts ...ContextOption) {
	for _, opt := range opts {
		opt.Apply(ctx)
	}
}

type ControllerFunc func(ctx *gin.Context) interface{}

func ControllerHandler(controllerFunc ControllerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		data := controllerFunc(ctx)

		switch r := data.(type) {
		case renders.ErrorRender:
			loggers.LogRequestErr(ctx, r.Err())
			r.Render(ctx)
		case renders.Render:
			r.Render(ctx)
		case ginRender.Render:
			ctx.Render(http.StatusOK, r)
		case error:
			er := response.UnknownError.WithErr(r)
			er.Render(ctx)
		default:
			cr := renders.JSON{Data: data}
			ctx.Render(http.StatusOK, cr)
		}
	}

}
