package tranx

import (
	"os"
	"fmt"
	"log"
	"sync"
	"time"
	"context"
	"net/http"
)

// Server is contains the necessary data to run the tranx server.
type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
	mailboxes  *Mailboxes
	ids        *IDs
	signal     chan os.Signal
}

// NewServer constructs a new Server struct and setups the routes.
func NewServer(port int) *Server {
	router := &http.ServeMux{}

	s := &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			Handler:      router,
		},
		router:    router,
		mailboxes: &Mailboxes{&sync.Map{}},
		ids:       &IDs{&sync.Map{}},
	}

	s.routes()

	return s
}

// Start runs the tranx server.
func (s *Server) Start() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-s.signal
		cancel()
	}()

	if err := serve(s, ctx); err != nil {
		log.Printf("Error serving Tran tranx server: %s\n", err)
	}
}

// serve is a helper function providing graceful shutdown of the server.
func serve(s *Server, ctx context.Context) (err error) {
	go func() {
		if err = s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Serving Tran: %s\n", err)
		}
	}()

	log.Printf("Tran Tranx Server started at \"%s\" \n", s.httpServer.Addr)
	<-ctx.Done()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = s.httpServer.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Tran tranx shutdown failed: %s", err)
	}

	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}
