package hub

import (
	"encoding/json"
	"github.com/Velocityofpie/chaudr/log"
	"github.com/gorilla/websocket"
)

const ErrorType = "error"
const MessageType = "message"

type WebsocketMessage struct {
	Type     string `json:"type"`
	Username string `json:"username,omitempty"`
	Message  string `json:"message"`
}

func CreateErrorMessage(message string) WebsocketMessage {
	return WebsocketMessage{
		Type:    ErrorType,
		Message: message,
	}
}

func SendErrorMessageAndClose(conn *websocket.Conn, errMessage string) {
	errWebsocketMsg := CreateErrorMessage(errMessage)
	errMsgJson, err := json.Marshal(errWebsocketMsg)
	if err != nil {
		log.Logger.Errorf("cannot marshal error message: %s", errMessage)
	}
	if err := conn.WriteMessage(websocket.TextMessage, errMsgJson); err != nil {
		log.Logger.Errorf("cannot write to client connection: %v", err)
	}
	if err := conn.Close(); err != nil {
		log.Logger.Errorf("cannot close client connection: %v", err)
	}
}

func CreateMessage(username, message string) WebsocketMessage {
	return WebsocketMessage{
		Type:     MessageType,
		Username: username,
		Message:  message,
	}
}
