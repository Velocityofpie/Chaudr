package hub

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RoomHub maintains the set of active clients and broadcasts messages to the
// clients.
type RoomHub struct {
	// Registered clients.
	// TODO: maybe need to control access to this? sync.Map?
	clients map[*ConnectedMember]bool

	// Inbound messages from the clients.
	broadcast chan WebsocketMessage

	// Register requests from the clients.
	register chan *ConnectedMember

	// Unregister requests from clients.
	unregister chan *ConnectedMember

	// Room Id
	roomId uint

	// mutex for controlling access to handshake
	m sync.Mutex

	// time when handshake was initiated
	handshakeStart time.Time

	// code for handshake. If code is 0, handshake is not initiated
	code uint32
}

func (h *RoomHub) ActiveClients() []string {
	clients := make([]string, 0, len(h.clients))
	for c, _ := range h.clients {
		clients = append(clients, c.member)
	}
	return clients
}

var HandshakeInProgress = errors.New("handshake already in progress")
var HandshakeCodeExpired = errors.New("handshake code expired")

func (h *RoomHub) SetHandshakeCode() (uint32, error) {
	h.m.Lock()
	defer h.m.Unlock()
	// check if handshake was initiated
	if !h.handshakeStart.IsZero() && time.Since(h.handshakeStart) < 30*time.Second {
		return 0, HandshakeInProgress
	}
	h.code = rand.Uint32()
	for h.code == 0 || h.code == math.MaxUint32 {
		h.code = rand.Uint32()
	}
	h.handshakeStart = time.Now()
	return h.code, nil
}

func (h *RoomHub) GetHandshakeCode() (uint32, error) {
	h.m.Lock()
	defer h.m.Unlock()
	if !h.handshakeStart.IsZero() && time.Since(h.handshakeStart) > 30*time.Second {
		h.code = 0
		return 0, HandshakeCodeExpired
	}
	return h.code, nil
}

// ResetHandshakeCode should be used only when handshake is completed early!!!!
func (h *RoomHub) ResetHandshakeCode() {
	h.m.Lock()
	defer h.m.Unlock()
	h.code = 0
	h.handshakeStart = time.Time{}
}

func NewHub(roomId uint) *RoomHub {
	return &RoomHub{
		broadcast:  make(chan WebsocketMessage),
		register:   make(chan *ConnectedMember),
		unregister: make(chan *ConnectedMember),
		clients:    make(map[*ConnectedMember]bool),
		roomId:     roomId,
	}
}

func (h *RoomHub) CheckUsernameTaken(username string) bool {
	for client := range h.clients {
		if client.member == username {
			return true
		}
	}
	return false
}

func (h *RoomHub) GetRoomId() uint {
	return h.roomId
}

func (h *RoomHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			// TODO: when no clients are left, then destroy the hub to save resources
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
				h.broadcast <- CreateMessage("hibot", "hi from hibot")
			}
		}
	}()
}
