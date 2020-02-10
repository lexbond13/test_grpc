package webwire

import logger2 "dancerate/logger"

var logger *logger2.Logger
var serverAddr string

func Init(serverAddress string, loggerInstance *logger2.Logger) {
	serverAddr = serverAddress
	logger = loggerInstance
}
