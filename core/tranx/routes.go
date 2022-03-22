package tranx

import "github.com/abdfnx/tran/tools"

func (s *Server) routes() {
	s.router.HandleFunc("/establish-sender", tools.WebsocketHandler(s.handleEstablishSender()))
	s.router.HandleFunc("/establish-receiver", tools.WebsocketHandler(s.handleEstablishReceiver()))
}
