package sender

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"net/http"
)

// handleTransfer creates a HandlerFunc to handle serving the transfer of files over a websocket connection
func (s *Sender) handleTransfer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.receiverIP.Equal(net.ParseIP(r.RemoteAddr)) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "No Tran for You!")
			log.Printf("Unauthorized Tran attempt from alien species with IP: %s\n", r.RemoteAddr)

			return
		}

		wsConn, err := s.senderServer.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Unable to initialize Tran due to technical error: %s\n", err)
			s.closeServer <- syscall.SIGTERM

			return
		}

		// Start transfer sequence.
		s.Transfer(wsConn)
	}
}
