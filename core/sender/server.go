package sender

import (
	"fmt"
	"log"
	"net"
	"time"
	"syscall"
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

// Server specifies the webserver that will be used for direct file transfer.
type Server struct {
	server   *http.Server
	router   *http.ServeMux
	upgrader websocket.Upgrader
}

// Specifies the necessary options for initializing the webserver.
type ServerOptions struct {
	port       int
	receiverIP net.IP
}

// Start starts the sender.Server webserver and setups graceful shutdown
func (s *Sender) StartServer() error {
	if s.senderServer == nil {
		return fmt.Errorf("start called with uninitialized senderServer")
	}

	// context used for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-s.closeServer
		cancel()
	}()

	// serve the webserver, and report errors
	if err := serve(s, ctx); err != nil {
		return err
	}

	return nil
}

func (s *Sender) CloseServer() {
	s.closeServer <- syscall.SIGTERM
}

// serve is helper function that serves the webserver while providing graceful shutdown.
func serve(s *Sender, ctx context.Context) (err error) {
	if s.senderServer == nil {
		return fmt.Errorf("serve called with uninitialized senderServer")
	}

	go func() {
		if err = s.senderServer.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Tran sender-server crashed due to an error: %s\n", err)
		}
	}()

	<-ctx.Done() // wait for the shutdown sequence to start.

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	// shutdown and report errors
	if err = s.senderServer.server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Tran shutdown sequence failed to due error: %s\n", err)
	}

	// strip error in this case, as we deal with this gracefully
	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}
