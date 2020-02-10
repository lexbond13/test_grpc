package server

import (
	"dancerate/api/handler"
	"github.com/fasthttp/router"
)

// GetRoutes
func GetRouter() *router.Router {

	appRouter := router.New()

	// Add routes here, may be use a groups
	appRouter.GET("/", handler.Index)
	appRouter.GET("/hello/:name", handler.Hello)

	return appRouter
}
