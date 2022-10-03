package common

import (
	"github.com/gin-gonic/gin"
)

func GetRouterPath(ctx *gin.Context) string {
	reqCtx := ctx.Request.Context()
	value := reqCtx.Value("router_path")
	routerPath, ok := value.(string)
	if !ok {
		return ""
	}
	return routerPath
}
