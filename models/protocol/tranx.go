package protocol

import (
	"net"

	"github.com/gorilla/websocket"
)

type TranxMessageType int

const (
	TranxToSenderBind TranxMessageType = iota // An ID for this connection is bound and communicated
	SenderToTranxEstablish   // Sender has generated and hashed password
	ReceiverToTranxEstablish // Passsword has been communicated to receiver who has hashed it
	TranxToSenderReady       // Tranx announces to sender that receiver is connected
	SenderToTranxPAKE        // Sender sends PAKE information to tranx
	TranxToReceiverPAKE      // Tranx forwards PAKE information to receiver
	ReceiverToTranxPAKE      // Receiver sends PAKE information to tranx
	TranxToSenderPAKE        // Tranx forwards PAKE information to receiver
	SenderToTranxSalt        // Sender sends cryptographic salt to tranx
	TranxToReceiverSalt      // Rendevoux forwards cryptographic salt to receiver
	ReceiverToTranxClose     // Receiver can connect directly to sender, close receiver connection -> close sender connection
	SenderToTranxClose       // Transit sequence is completed, close sender connection -> close receiver connection
)

type TranxMessage struct {
	Type    TranxMessageType `json:"type"`
	Payload interface{}      `json:"payload"`
}

type TranxClient struct {
	Conn *websocket.Conn
	IP   net.IP
}

type TranxSender struct {
	TranxClient
	Port int
}

type TranxReceiver = TranxClient

/* [💻 Receiver <-> Sender 💻] messages */

type PasswordPayload struct {
	Password string `json:"password"`
}
type PakePayload struct {
	Bytes []byte `json:"pake_bytes"`
}

type SaltPayload struct {
	Salt []byte `json:"salt"`
}

/* [Tranx -> Sender] messages */

type TranxToSenderBindPayload struct {
	ID int `json:"id"`
}
