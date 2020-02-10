package webwire

import "context"

// handleSignal handles incoming signals
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSignal(clt *Client, msg *Message) {
	srv.opsLock.Lock()
	// Ignore incoming signals during shutdown
	if srv.shutdown {
		srv.opsLock.Unlock()
		return
	}
	srv.currentOps++
	srv.opsLock.Unlock()

	srv.impl.OnSignal(
		context.Background(),
		clt,
		msg,
	)

	// Mark signal as done and shutdown the server if scheduled and no ops are left
	srv.opsLock.Lock()
	srv.currentOps--
	if srv.shutdown && srv.currentOps < 1 {
		close(srv.shutdownRdy)
	}
	srv.opsLock.Unlock()
}
