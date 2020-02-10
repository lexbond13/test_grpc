package webwire

import (
	"context"
	"fmt"
	wwr "github.com/qbeon/webwire-go"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type WWServer struct{}

// OnOptions implements the webwire.ServerImplementation interface.
// Sets HTTP access control headers to satisfy CORS
func (srv *WWServer) OnOptions(resp http.ResponseWriter) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "WEBWIRE")
}

// OnSignal implements the webwire.ServerImplementation interface
// Does nothing, not needed in this example
func (srv *WWServer) OnSignal(
	_ context.Context,
	_ *wwr.Client,
	_ *wwr.Message,
) {
}

// OnClientConnected implements the webwire.ServerImplementation interface.
// Does nothing, not needed in this example
func (srv *WWServer) OnClientConnected(client *wwr.Client) {}

// OnClientDisconnected implements the webwire.ServerImplementation interface
// Does nothing, not needed in this example
func (srv *WWServer) OnClientDisconnected(client *wwr.Client) {}

// BeforeUpgrade implements the webwire.ServerImplementation interface.
// Must return true to ensure incoming connections are accepted
func (srv *WWServer) BeforeUpgrade(resp http.ResponseWriter, req *http.Request) bool {
	return true
}

// OnRequest implements the webwire.ServerImplementation interface.
// Returns the received message back to the client
func (srv *WWServer) OnRequest(_ context.Context,client *wwr.Client,message *wwr.Message) (response wwr.Payload, err error) {
	logger.Info().Msg(fmt.Sprintf("Replied to client: %s, message: %s", client.Info().RemoteAddr, string(message.Payload.Data)))

	// Reply to the request using the same data and encoding
	return message.Payload, nil
}


func RunServer() {

	// Setup a new webwire server instance
	server, err := wwr.NewServer(
		&WWServer{},
		wwr.ServerOptions{
			Address: serverAddr,
		},
	)
	if err != nil {
		panic(fmt.Errorf("Failed setting up WebWire server: %s", err))
	}

	// Listen for OS signals and shutdown server in case of demanded termination
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-osSignals
		logger.Info().Msg(fmt.Sprintf("Termination demanded by the OS (%s), shutting down...", sig))
		if err := server.Shutdown(); err != nil {
			logger.Info().Msg(fmt.Sprintf("Error during server shutdown: %s", err))
		}
		logger.Info().Msg(fmt.Sprintf("Server gracefully terminated"))
	}()

	// Launch echo server
	logger.Info().Msg(fmt.Sprintf("Listening on %s", server.Addr().String()))
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
