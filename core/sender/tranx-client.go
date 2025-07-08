package sender

import (
	"fmt"
	"net"

	"github.com/schollz/pake/v3"
	"github.com/gorilla/websocket"
	"github.com/abdfnx/tran/tools"
	"github.com/abdfnx/tran/models"
	"github.com/abdfnx/tran/core/crypt"
	"github.com/abdfnx/tran/models/protocol"
)

// ConnectToTranx, establishes the connection with the tranx server.
// Parameters:
// tranxAddress 	-   IP or hostname of the tranx server
// tranxPort 		- 	port of the tranx server
// startServerCh    -   channel to communicate to the caller when to start the server, and with which options.
// passwordCh       -   channel to communicate the password to the caller.
// startServerCh    -   channel to communicate to the caller when to start the server, and with which options.
// payloadReady    	-   channel over which the caller can communicate when the payload is ready.
// relayCh         	-   channel to communicated if we are using relay (tranx) for transfer.
func (s *Sender) ConnectToTranx(
	tranxAddress string,
	tranxPort int,
	passwordCh chan<- models.Password,
	startServerCh chan<- ServerOptions,
	payloadReady <-chan bool,
	relayCh chan<- *websocket.Conn,
) error {
	// establish websocket connection to tranx server
	wsConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/establish-sender", tranxAddress, tranxPort), nil)
	if err != nil {
		return err
	}

	// bind connection
	tranxMsg, err := tools.ReadTranxMessage(wsConn, protocol.TranxToSenderBind)
	if err != nil {
		return err
	}

	bindPayload := protocol.TranxToSenderBindPayload{}
	err = tools.DecodePayload(tranxMsg.Payload, &bindPayload)
	if err != nil {
		return err
	}

	// establish sender
	password := tools.GeneratePassword(bindPayload.ID)
	hashed := tools.HashPassword(password)

	wsConn.WriteJSON(protocol.TranxMessage{
		Type: protocol.SenderToTranxEstablish,
		Payload: protocol.PasswordPayload{
			Password: hashed,
		},
	})

	// send the generated password to the UI so it can be displayed
	passwordCh <- password

	// setup the encryption
	err = s.establishSecureConnection(wsConn, password)
	if err != nil {
		return err
	}

	// do the transfer handshake over the tranx
	err = s.doHandshake(wsConn, payloadReady, startServerCh)
	if err != nil {
		return err
	}

	transferMsg, err := tools.ReadEncryptedMessage(wsConn, s.crypt)
	if err != nil {
		return err
	}

	switch transferMsg.Type {
	// we will do direct communication with the receiver
		case protocol.ReceiverDirectCommunication:
			close(relayCh)
			tools.WriteEncryptedMessage(wsConn, protocol.TransferMessage{Type: protocol.SenderDirectAck}, s.crypt)

			return nil

		// we will do relay communication with receiver using the same websocket connection as with tranx
		case protocol.ReceiverRelayCommunication:
			tools.WriteEncryptedMessage(wsConn, protocol.TransferMessage{Type: protocol.SenderRelayAck}, s.crypt)
			relayCh <- wsConn

			return nil

		default:
			return protocol.NewWrongMessageTypeError(
				[]protocol.TransferMessageType{protocol.ReceiverDirectCommunication, protocol.ReceiverRelayCommunication},
				transferMsg.Type)
	}
}

// establishSecureConnection setups the PAKE2 key exchange and the crypt struct in the sender.
func (s *Sender) establishSecureConnection(wsConn *websocket.Conn, password models.Password) error {
	// init PAKE2 (NOTE: This takes a couple of seconds, here it is fine as we have to wait for the receiver)
	pake, err := pake.InitCurve([]byte(password), 0, "p256")

	if err != nil {
		return err
	}

	// Wait for receiver to be ready to exchange crypto information.
	msg, err := tools.ReadTranxMessage(wsConn, protocol.TranxToSenderReady)
	if err != nil {
		return err
	}

	// PAKE sender -> receiver.
	wsConn.WriteJSON(protocol.TranxMessage{
		Type: protocol.SenderToTranxPAKE,
		Payload: protocol.PakePayload{
			Bytes: pake.Bytes(),
		},
	})

	// PAKE receiver -> sender.
	msg, err = tools.ReadTranxMessage(wsConn, protocol.TranxToSenderPAKE)
	if err != nil {
		return err
	}

	pakePayload := protocol.PakePayload{}
	err = tools.DecodePayload(msg.Payload, &pakePayload)
	if err != nil {
		return err
	}

	err = pake.Update(pakePayload.Bytes)
	if err != nil {
		return err
	}

	// Setup crypt.Crypt struct in Sender.
	sessionkey, err := pake.SessionKey()
	if err != nil {
		return err
	}

	s.crypt, err = crypt.New(sessionkey)
	if err != nil {
		return err
	}

	// Send salt to receiver.
	wsConn.WriteJSON(protocol.TranxMessage{
		Type: protocol.SenderToTranxSalt,
		Payload: protocol.SaltPayload{
			Salt: s.crypt.Salt,
		},
	})

	return nil
}

// doHandshake does the transfer handshake over the tranx connection
func (s *Sender) doHandshake(wsConn *websocket.Conn, payloadReady <-chan bool, startServerCh chan<- ServerOptions) error {
	transferMsg, err := tools.ReadEncryptedMessage(wsConn, s.crypt)
	if err != nil {
		return err
	}

	if transferMsg.Type != protocol.ReceiverHandshake {
		return protocol.NewWrongMessageTypeError([]protocol.TransferMessageType{protocol.ReceiverHandshake}, transferMsg.Type)
	}

	handshakePayload := protocol.ReceiverHandshakePayload{}
	err = tools.DecodePayload(transferMsg.Payload, &handshakePayload)
	if err != nil {
		return err
	}

	senderPort, err := tools.GetOpenPort()
	if err != nil {
		return err
	}

	// wait for payload to be ready
	<-payloadReady
	startServerCh <- ServerOptions{port: senderPort, receiverIP: handshakePayload.IP}

	tcpAddr, _ := wsConn.LocalAddr().(*net.TCPAddr)
	handshake := protocol.TransferMessage{
		Type: protocol.SenderHandshake,
		Payload: protocol.SenderHandshakePayload{
			IP:          tcpAddr.IP,
			Port:        senderPort,
			PayloadSize: s.payloadSize,
		},
	}

	tools.WriteEncryptedMessage(wsConn, handshake, s.crypt)

	return nil
}
