package hub

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RoomHub maintains the set of active clients and broadcasts messages to the
// clients.
type RoomHub struct {
	// Registered clients.
	clients map[*ConnectedMember]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *ConnectedMember

	// Unregister requests from clients.
	unregister chan *ConnectedMember
}

func NewHub() *RoomHub {
	return &RoomHub{
		broadcast:  make(chan []byte),
		register:   make(chan *ConnectedMember),
		unregister: make(chan *ConnectedMember),
		clients:    make(map[*ConnectedMember]bool),
	}
}

func (h *RoomHub) CheckUsernameTaken(username string) bool {
	for client := range h.clients {
		if client.member.Username == username {
			return true
		}
	}
	return false
}

func (h *RoomHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// BotHandler handles websocket requests from the peer.
func BotHandler(hubMap *sync.Map, w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	username := r.URL.Query().Get("username")

	if username == "" {
		username = "hibot"
	}

	if roomId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("join request has blank room id"))
		return
	}

	id, err := strconv.ParseUint(roomId, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not parse room id: %s", roomId)))
		return
	}

	hub, ok := hubMap.Load(uint(id))
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unknown room id %d", id)))
		return
	}

	h := hub.(*RoomHub)

	if h.CheckUsernameTaken(username) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("member is already connected: %s", username)))
		return
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				h.broadcast <- []byte("hi from hibot")
			}
		}
	}()
}
