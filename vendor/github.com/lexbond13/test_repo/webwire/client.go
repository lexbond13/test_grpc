package webwire

import (
	"fmt"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

var wwClient *WWClient

// EchoClient implements the wwrclt.Implementation interface
type WWClient struct {
	connection *wwrclt.Client
}

// NewEchoClient constructs and returns a new echo client instance
func NewWWClient() *WWClient {

	wwClient = &WWClient{}
	wwClient.connection = wwrclt.NewClient(
		serverAddr,
		wwClient,
		wwrclt.Options{
			DefaultRequestTimeout: 10 * time.Second,
			ReconnectionInterval: 2 * time.Second,
		},
	)

	return wwClient
}

func GetClient() *WWClient {
	return wwClient
}

// OnDisconnected implements the wwrclt.Implementation interface
func (clt *WWClient) OnDisconnected() {}

// OnSessionClosed implements the wwrclt.Implementation interface
func (clt *WWClient) OnSessionClosed() {}

// OnSessionCreated implements the wwrclt.Implementation interface
func (clt *WWClient) OnSessionCreated(_ *wwr.Session) {}

// OnSignal implements the wwrclt.Implementation interface
func (clt *WWClient) OnSignal(_ wwr.Payload) {}

// Request sends a message to the server and returns the reply.
// panics if the request fails for whatever reason
func (clt *WWClient) Request(message string) wwr.Payload {
	// Define a payload to be sent to the server, use default binary encoding
	payload := wwr.Payload{
		Data: []byte(message),
	}

	logger.Info().Msg(fmt.Sprintf(
		"Sent request:   '%s' (%d)",
		string(payload.Data),
		len(payload.Data),
	))

	// Send request and await reply
	reply, err := clt.connection.Request("", payload)
	if err != nil {
		panic(fmt.Errorf("Request failed: %s", err))
	}

	return reply
}
