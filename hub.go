package main

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

func newHub() *RoomHub {
	return &RoomHub{
		broadcast:  make(chan []byte),
		register:   make(chan *ConnectedMember),
		unregister: make(chan *ConnectedMember),
		clients:    make(map[*ConnectedMember]bool),
	}
}

func (h *RoomHub) run() {
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
