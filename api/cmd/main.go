package main

import (
	"github.com/lexbond13/test_repo/api/server"
	lr "github.com/lexbond13/test_repo/logger"
	"github.com/lexbond13/test_repo/webwire"
)

func main() {

	logger := lr.New(true) // TODO replace this flag to config
	webwire.Init(":8091", logger) // TODO replace port to config

	//go webwire.RunServer()
	//wwClient := webwire.NewWWClient()
	//wwClient.Request("this is message for websockets")

	logger.Info().Msg("Starting listen on port :8090") // TODO set port to config

	if err := server.Run(); err != nil {
		logger.Error().Err(err)
		panic(err)
	}
}
