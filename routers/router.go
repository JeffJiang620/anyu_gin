package routers

import (
	"context"
	"path"
	"strings"
	"unsafe"

	"github.com/anyufly/gin_common/common"
	"github.com/anyufly/gin_common/controllers"
	"github.com/anyufly/gin_common/middlewares"
	"github.com/gin-gonic/gin"
)

type RouteDesc struct {
	Method     string
	MiddleWare []interface{}
	Controller []controllers.ControllerFunc
}

type Router interface {
	GroupName() string
	GroupConfig() map[string][]RouteDesc
	GroupMiddleware() []middlewares.IMiddleWare
}

func handleRouter(version *gin.RouterGroup, router Router) {
	groupName := router.GroupName()
	groupConfigs := router.GroupConfig()
	groupMiddlewares := router.GroupMiddleware()
	group := handleGroupName(version, groupName)
	middlewareFlag := handleGroupMiddleware(group, groupMiddlewares...)
	handleGroupConfigs(group, groupConfigs, middlewareFlag)
}

func handleGroupMiddleware(group *gin.RouterGroup, midList ...middlewares.IMiddleWare) map[int]bool {
	var middlewareHandlers []gin.HandlerFunc
	var middlewareFlag = make(map[int]bool)
	for _, middleware := range midList {
		mp := (*int)(unsafe.Pointer(&middleware))
		if _, ok := middlewareFlag[*mp]; !ok {
			middlewareFlag[*mp] = true
			middlewareHandlers = append(middlewareHandlers, middlewares.MiddlewareHandler(middleware))
		}
	}
	group.Use(middlewareHandlers...)
	return middlewareFlag
}

func handleGroupName(version *gin.RouterGroup, groupName string) *gin.RouterGroup {
	trimName := strings.Trim(groupName, " ")
	if trimName == "" {
		return version
	}
	return version.Group(trimName)
}

type RequestFunc func(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

func addRouterFullPath(basePath string, relativePath string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqCtx := ctx.Request.Context()
		valCtx := context.WithValue(reqCtx, common.ContextKey("router_path"), calculateAbsolutePath(basePath, relativePath))
		ctx.Request = ctx.Request.WithContext(valCtx)
	}
}

func calculateAbsolutePath(basePath string, relativePath string) string {
	return joinPaths(basePath, relativePath)
}

func joinPaths(basePath, relativePath string) string {
	if relativePath == "" {
		return basePath
	}

	finalPath := path.Join(basePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func getFuncByMethod(method string) RequestFunc {
	upperMethod := strings.ToUpper(method)

	return func(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
		return group.Handle(upperMethod, relativePath, handlers...)
	}
}

func handleGroupConfigs(group *gin.RouterGroup, groupConfigs map[string][]RouteDesc, alreadyUseMiddlewareFlag map[int]bool) {
	for relativePath, routeDescribes := range groupConfigs {
		for _, routeDesc := range routeDescribes {
			method := routeDesc.Method

			var controllerHandlers []gin.HandlerFunc
			var middlewareFlag = make(map[int]bool)

			for _, middleware := range routeDesc.MiddleWare {
				mp := (*int)(unsafe.Pointer(&middleware))
				if _, ok := middlewareFlag[*mp]; !ok {
					if _, alreadyUse := alreadyUseMiddlewareFlag[*mp]; !alreadyUse {
						middlewareFlag[*mp] = true
						switch middleware.(type) {
						case middlewares.IMiddleWare:
							controllerHandlers = append(controllerHandlers,
								middlewares.MiddlewareHandler(middleware.(middlewares.IMiddleWare)))
						case middlewares.MiddlewareFunc:
							fn := middleware.(middlewares.MiddlewareFunc)
							mid := fn()
							controllerHandlers = append(controllerHandlers, middlewares.MiddlewareHandler(mid))
						default:

						}

					}
				}
			}

			var controllerFlag = make(map[int]bool)
			for _, controller := range routeDesc.Controller {
				cp := (*int)(unsafe.Pointer(&controller))
				if _, ok := controllerFlag[*cp]; !ok {
					controllerFlag[*cp] = true
					controllerHandlers = append(controllerHandlers, controllers.ControllerHandler(controller))
				}
			}

			process := getFuncByMethod(method)
			newHandlers := make([]gin.HandlerFunc, 0, len(controllerHandlers)+1)
			newHandlers = append(newHandlers, addRouterFullPath(group.BasePath(), relativePath))
			newHandlers = append(newHandlers, controllerHandlers...)

			process(group, relativePath, newHandlers...)
		}
	}
}

func CombineRouters(engine *gin.Engine, basePath string, routers ...Router) {
	v1 := engine.Group(basePath)
	for _, router := range routers {
		handleRouter(v1, router)
	}
}
