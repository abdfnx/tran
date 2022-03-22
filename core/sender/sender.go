package sender

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/scmn-dev/tran/core/crypt"
	"github.com/scmn-dev/tran/models"
)

// Sender represents the sender client, handles tranx communication and file transfer.
type Sender struct {
	payload      io.Reader
	payloadSize  int64
	senderServer *Server
	closeServer  chan os.Signal
	receiverIP   net.IP
	tranxAddress string
	tranxPort    int
	ui           chan<- UIUpdate
	crypt        *crypt.Crypt
	state        TransferState
}

// NewSender returns a bare bones Sender.
func NewSender(programOptions models.TranOptions) *Sender {
	closeServerCh := make(chan os.Signal, 1)
	signal.Notify(closeServerCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return &Sender{
		closeServer:       closeServerCh,
		tranxAddress: programOptions.TranxAddress,
		tranxPort:    programOptions.TranxPort,
		state:             Initial,
	}
}

// WithPayload specifies the payload that will be transfered.
func WithPayload(s *Sender, payload io.Reader, payloadSize int64) *Sender {
	s.payload = payload
	s.payloadSize = payloadSize

	return s
}

// WithServer specifies the option to run the sender by hosting a server which the receiver establishes a connection to.
func WithServer(s *Sender, options ServerOptions) *Sender {
	s.receiverIP = options.receiverIP
	router := &http.ServeMux{}
	s.senderServer = &Server{
		router: router,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", options.port),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			Handler:      router,
		},
		upgrader: websocket.Upgrader{},
	}

	// setup routes
	router.HandleFunc("/tran", s.handleTransfer())
	return s
}

// WithUI specifies the option to run the sender with an UI channel that reports the state of the transfer.
func WithUI(s *Sender, ui chan<- UIUpdate) *Sender {
	s.ui = ui

	return s
}

func (s *Sender) TranxAddress() string {
	return s.tranxAddress
}

func (s *Sender) TranxPort() int {
	return s.tranxPort
}

// updateUI is a helper function that checks if we have a UI channel and reports the state.
func (s *Sender) updateUI(progress ...float32) {
	if s.ui == nil {
		return
	}

	var p float32

	if len(progress) > 0 {
		p = progress[0]
	}

	s.ui <- UIUpdate{State: s.state, Progress: p}
}
