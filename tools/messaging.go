package tools

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/scmn-dev/tran/core/crypt"
	"github.com/scmn-dev/tran/models/protocol"
)

func ReadTranxMessage(wsConn *websocket.Conn, expected protocol.TranxMessageType) (protocol.TranxMessage, error) {
	msg := protocol.TranxMessage{}
	err := wsConn.ReadJSON(&msg)

	if err != nil {
		return protocol.TranxMessage{}, err
	}

	if msg.Type != expected {
		return protocol.TranxMessage{}, fmt.Errorf("expected message type: %d. Got type: %d", expected, msg.Type)
	}

	return msg, nil
}

func WriteEncryptedMessage(wsConn *websocket.Conn, msg protocol.TransferMessage, crypt *crypt.Crypt) error {
	json, err := json.Marshal(msg)

	if err != nil {
		return nil
	}

	enc, err := crypt.Encrypt(json)
	if err != nil {
		return err
	}

	wsConn.WriteMessage(websocket.BinaryMessage, enc)

	return nil
}

func ReadEncryptedMessage(wsConn *websocket.Conn, crypt *crypt.Crypt) (protocol.TransferMessage, error) {
	_, enc, err := wsConn.ReadMessage()

	if err != nil {
		return protocol.TransferMessage{}, err
	}

	dec, err := crypt.Decrypt(enc)
	if err != nil {
		return protocol.TransferMessage{}, err
	}

	msg := protocol.TransferMessage{}
	err = json.Unmarshal(dec, &msg)

	if err != nil {
		return protocol.TransferMessage{}, err
	}

	return msg, nil
}
