package server

import (
	"github.com/valyala/fasthttp"
)

// Run starting listen server
func Run() error {

	appRouter := GetRouter()

	err := fasthttp.ListenAndServe(":8090", appRouter.Handler)
	if err != nil {
		return err
	}

	return nil
}
