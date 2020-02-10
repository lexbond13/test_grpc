package webwire

import logger2 "test_repo/logger"

var logger *logger2.Logger
var serverAddr string

func Init(serverAddress string, loggerInstance *logger2.Logger) {
	serverAddr = serverAddress
	logger = loggerInstance
}
