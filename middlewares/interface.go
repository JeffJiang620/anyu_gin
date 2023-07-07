package middlewares

import "github.com/gin-gonic/gin"

type IMiddleWare interface {
	Before(ctx *gin.Context) interface{}
	After(ctx *gin.Context) interface{}
	DeniedBeforeAbortContext() bool
	AllowAfterAbortContext() bool
}

type MiddlewareFunc func() IMiddleWare
