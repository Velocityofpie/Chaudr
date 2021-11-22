// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hub

import (
	"encoding/json"
	"github.com/Velocityofpie/chaudr/config"
	"github.com/Velocityofpie/chaudr/log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func init() {
	if config.DebugMode {
		Upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}
}

// ConnectedMember is a middleman between the websocket connection and the hub.
type ConnectedMember struct {
	hub *RoomHub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan WebsocketMessage

	member string

	registered bool
}

func NewConnectedMember(hub *RoomHub, conn *websocket.Conn, send chan WebsocketMessage, member string) ConnectedMember {
	return ConnectedMember{
		hub:    hub,
		conn:   conn,
		send:   send,
		member: member,
	}
}

func (c *ConnectedMember) Register() {
	if !c.registered {
		c.hub.register <- c
	}
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go c.writePump()
	go c.readPump()
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *ConnectedMember) readPump() {
	defer func() {
		c.hub.unregister <- c
		if err := c.conn.Close(); err != nil {
			log.Logger.Errorf("could not close connection to client: %v", err)
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msgType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Logger.Errorf("error: %v", err)
			}
			break
		}
		if msgType != websocket.TextMessage {
			errMsg := CreateErrorMessage("server only supports text message")
			errMsgJson, err := json.Marshal(errMsg)
			if err != nil {
				log.Logger.Errorf("cannot marshal error message: %v", errMsg)
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, errMsgJson); err != nil {
				log.Logger.Errorf("cannot write to client connection: %v", err)
			}
			if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
				log.Logger.Errorf("cannot write to client connection: %v", err)
			}
			break
		}
		var clientMsg WebsocketMessage
		if err := json.Unmarshal(message, &clientMsg); err != nil {
			log.Logger.Errorf("cannot marshal client send message: %s", message)
			errMsg := CreateErrorMessage("malformed message sent")
			errMsgJson, err := json.Marshal(errMsg)
			if err != nil {
				log.Logger.Errorf("cannot marshal error message: %v", errMsg)
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, errMsgJson); err != nil {
				log.Logger.Errorf("cannot write to client connection: %v", err)
			}
			if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
				log.Logger.Errorf("cannot write to client connection: %v", err)
			}
			break
		}
		log.Logger.Debugf("client sent message: %v", clientMsg)
		c.hub.broadcast <- clientMsg
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *ConnectedMember) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			log.Logger.Errorf("could not close connection to client: %v", err)
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			msgJson, err := json.Marshal(message)
			if err != nil {
				log.Logger.Errorf("could not marshal message: %v", message)
			} else {
				w.Write(msgJson)
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				msgJson, err := json.Marshal(<-c.send)
				if err != nil {
					log.Logger.Errorf("could not marshal message: %v", message)
				} else {
					w.Write(msgJson)
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
