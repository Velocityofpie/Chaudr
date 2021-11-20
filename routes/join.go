package routes

import (
	"fmt"
	"github.com/Velocityofpie/chaudr/hub"
	"github.com/Velocityofpie/chaudr/repository"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// joinRoomHandler handles websocket requests from the peer.
func joinRoomHandler(hubMap *sync.Map, w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	username := r.URL.Query().Get("username")

	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("join request has blank username"))
		return
	}

	if roomId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("join request has blank room id"))
		return
	}

	// TODO: check if the given room id and username pair is valid
	// ...

	id, err := strconv.ParseUint(roomId, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not parse room id: %s", roomId)))
		return
	}

	// TODO: if room hub does not exist, create it. Implement this after repository
	interfaceHub, ok := hubMap.Load(uint(id))
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unknown room id %d", id)))
		return
	}

	h := interfaceHub.(*hub.RoomHub)

	if h.CheckUsernameTaken(username) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("member is already connected: %s", username)))
		return
	}

	conn, err := hub.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := hub.NewConnectedMember(
		h,
		conn,
		make(chan []byte, 256),
		repository.memberSqlModel{
			RoomID:   uint(id),
			Username: username,
		},
	)
	client.Register()
}
