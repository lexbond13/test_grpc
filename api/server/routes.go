package server

import (
	"github.com/fasthttp/router"
	"test_repo/api/handler"
)

// GetRoutes
func GetRouter() *router.Router {

	appRouter := router.New()

	// Add routes here, may be use a groups
	appRouter.GET("/", handler.Index)
	appRouter.GET("/hello/:name", handler.Hello)

	return appRouter
}
