package webwire

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

const protocolVersion = "1.4"

// server represents a headless WebWire server instance,
// where headless means there's no HTTP server that's hosting it
type server struct {
	impl              ServerImplementation
	httpServer        *http.Server
	listener          net.Listener
	sessionManager    SessionManager
	sessionKeyGen     SessionKeyGenerator
	sessionInfoParser SessionInfoParser

	// State
	addr            net.Addr
	shutdown        bool
	shutdownRdy     chan bool
	currentOps      uint32
	opsLock         sync.Mutex
	clientsLock     *sync.Mutex
	clients         []*Client
	sessionsEnabled bool
	sessionRegistry *sessionRegistry

	// Internals
	connUpgrader ConnUpgrader
	warnLog      *log.Logger
	errorLog     *log.Logger
}

func (srv *server) shutdownHTTPServer() error {
	if srv.httpServer == nil {
		return nil
	}
	if err := srv.httpServer.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("Couldn't properly shutdown HTTP server: %s", err)
	}
	return nil
}

// Run implements the Server interface
func (srv *server) Run() error {
	// Launch HTTP server
	if err := srv.httpServer.Serve(
		tcpKeepAliveListener{srv.listener.(*net.TCPListener)},
	); err != http.ErrServerClosed {
		return fmt.Errorf("HTTP Server failure: %s", err)
	}

	return nil
}

// Addr implements the Server interface
func (srv *server) Addr() net.Addr {
	return srv.addr
}

// Shutdown implements the Server interface
func (srv *server) Shutdown() error {
	srv.opsLock.Lock()
	srv.shutdown = true
	// Don't block if there's no currently processed operations
	if srv.currentOps < 1 {
		srv.opsLock.Unlock()
		return srv.shutdownHTTPServer()
	}
	srv.opsLock.Unlock()
	<-srv.shutdownRdy

	return srv.shutdownHTTPServer()
}

// ActiveSessionsNum implements the Server interface
func (srv *server) ActiveSessionsNum() int {
	return srv.sessionRegistry.activeSessionsNum()
}

// SessionConnectionsNum implements the Server interface
func (srv *server) SessionConnectionsNum(sessionKey string) int {
	return srv.sessionRegistry.sessionConnectionsNum(sessionKey)
}

// SessionConnections implements the Server interface
func (srv *server) SessionConnections(sessionKey string) []*Client {
	agents := srv.sessionRegistry.sessionConnections(sessionKey)
	if agents == nil {
		return nil
	}
	list := make([]*Client, len(agents))
	copy(list, agents)
	return list
}

// CloseSession implements the Server interface
func (srv *server) CloseSession(sessionKey string) int {
	agents := srv.sessionRegistry.sessionConnections(sessionKey)
	if agents == nil {
		return -1
	}
	for _, agent := range agents {
		agent.Close()
	}
	return len(agents)
}
