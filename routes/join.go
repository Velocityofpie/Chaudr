package routes

import (
	"encoding/json"
	"fmt"
	"github.com/Velocityofpie/chaudr/hub"
	"github.com/Velocityofpie/chaudr/log"
	"github.com/Velocityofpie/chaudr/repository"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
)

// joinRoomHandler handles websocket requests from the peer.
func joinRoomHandler(repo repository.RoomRepository, hubMap *sync.Map, w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	username := r.URL.Query().Get("username")

	conn, err := hub.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Logger.Errorf("could not upgrade connection: %s", err)
		return
	}

	if username == "" {
		hub.SendErrorMessageAndClose(conn, "join request has blank username")
		return
	}

	if roomId == "" {
		hub.SendErrorMessageAndClose(conn, "join request has blank room id")
		return
	}

	id, err := strconv.ParseUint(roomId, 10, 32)
	if err != nil {
		hub.SendErrorMessageAndClose(conn, fmt.Sprintf("could not parse room id: %s", roomId))
		return
	}

	// check if the given room id and username pair is valid
	room, err := repo.GetRoom(repository.Room{
		Id: uint(id),
	})
	if err != nil {
		hub.SendErrorMessageAndClose(conn, fmt.Sprintf("could not find room %s", roomId))
		return
	}
	var found bool
	for _, m := range room.Members {
		if m == username {
			found = true
			break
		}
	}
	if !found {
		hub.SendErrorMessageAndClose(conn, fmt.Sprintf("user %s is not part of the room %s", username, roomId))
		return
	}

	var h *hub.RoomHub

	interfaceHub, ok := hubMap.Load(uint(id))
	if !ok {
		newH := hub.NewHub(uint(id))
		hubMap.Store(uint(id), newH)
		h = newH
	} else {
		h = interfaceHub.(*hub.RoomHub)
	}

	if h.CheckUsernameTaken(username) {
		e := fmt.Sprintf("user is already connected to room %d", id)
		errMsg := hub.CreateErrorMessage(e)
		errMsgJson, err := json.Marshal(errMsg)
		if err != nil {
			log.Logger.Errorf("cannot marshal error message: %s", e)
		}
		if err := conn.WriteMessage(websocket.TextMessage, errMsgJson); err != nil {
			log.Logger.Errorf("cannot write to client connection: %v", err)
		}
		if err := conn.Close(); err != nil {
			log.Logger.Errorf("cannot close client connection: %v", err)
		}
		return
	}

	client := hub.NewConnectedMember(
		h,
		conn,
		make(chan hub.WebsocketMessage, 256),
		username,
	)
	client.Register()
}
