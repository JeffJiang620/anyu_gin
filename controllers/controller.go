package controllers

import (
	"github.com/anyufly/gin_common/apierr"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/anyufly/gin_common/loggers"
	"github.com/anyufly/gin_common/renders"
	"github.com/anyufly/gin_common/response"
	"github.com/gin-gonic/gin"
	ginRender "github.com/gin-gonic/gin/render"
)

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
		case validator.ValidationErrors:
			er := response.ParameterError.WithErr(r)
			er.Render(ctx)
		case *apierr.APIError:
			er := &response.ErrorResponse{}
			er.WithStatusCode(http.StatusOK).WithErr(r).Render(ctx)
		case error:
			loggers.LogRequestErr(ctx, r)
			er := response.UnknownError.WithErr(r)
			er.Render(ctx)
		default:
			cr := renders.JSON{Data: data}
			ctx.Render(http.StatusOK, cr)
		}
	}

}
