package routes

import (
	"encoding/json"
	"errors"
	"github.com/Velocityofpie/chaudr/hub"
	"github.com/Velocityofpie/chaudr/log"
	"github.com/Velocityofpie/chaudr/repository"
	"io"
	"net/http"
	"sync"
)

type addMemberRequestBody struct {
	Code     uint32 `json:"code"`
	Username string `json:"username"`
	RoomId   uint   `json:"roomId"`
}

// PUT /room/user
func addUserToRoomHandler(repo repository.RoomRepository, hubMap *sync.Map, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unexpected http method"))
		return
	}

	var body addMemberRequestBody
	contents, err := io.ReadAll(r.Body)
	if err != nil {
		log.Logger.Errorf("could not read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not read request body"))
		return
	}

	if err := json.Unmarshal(contents, &body); err != nil {
		log.Logger.Errorf("could not read response body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not read request body"))
		return
	}

	log.Logger.Debugf("join handshake request body: %v", body)

	if body.RoomId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body gave invalid room id"))
		return
	}

	if body.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body gave blank username"))
		return
	}

	if body.Code == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body gave invalid code"))
		return
	}

	room, err := repo.GetRoom(repository.Room{Id: body.RoomId})
	if err != nil {
		if errors.Is(err, repository.UnknownRoomId) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("room does not exist"))
		} else {
			log.Logger.Errorf("could not get room: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not get room"))
		}
		return
	}

	// check if user is already part of room
	for _, m := range room.Members {
		if m == body.Username {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("username is already in use"))
			return
		}
	}

	// check code
	var h *hub.RoomHub

	interfaceHub, ok := hubMap.Load(body.RoomId)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("room hub has not been loaded"))
		return
	} else {
		h = interfaceHub.(*hub.RoomHub)
	}

	code, err := h.GetHandshakeCode()

	log.Logger.Debugf("hub code is %d", code)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("handshake code is expired"))
		return
	}

	if code == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("handshake not initiated"))
		return
	}

	if body.Code != code {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("code does not match"))
		return
	}

	_, err = repo.AddMemberToRoom(room, body.Username)
	if err != nil {
		log.Logger.Errorf("could not add member %s to room %d", body.Username, room.Id)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("room hub has not been loaded"))
		return
	}

	w.Write([]byte("member added"))
	h.ResetHandshakeCode()
}

type startHandshakeRequestBody struct {
	Username string `json:"username"`
	RoomId   uint   `json:"roomId"`
}

type startHandshakeResponse struct {
	Code uint32 `json:"code"`
}

// PUT /room/handshake
func startHandshakeHandler(repo repository.RoomRepository, hubMap *sync.Map, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unexpected http method"))
		return
	}

	var body startHandshakeRequestBody
	contents, err := io.ReadAll(r.Body)
	if err != nil {
		log.Logger.Errorf("could not read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not read request body"))
		return
	}

	if err := json.Unmarshal(contents, &body); err != nil {
		log.Logger.Errorf("could not read response body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not read request body"))
		return
	}

	if body.RoomId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body gave invalid room id"))
		return
	}

	if body.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body gave blank username"))
		return
	}

	// check that user is part of room
	room, err := repo.GetRoom(repository.Room{Id: body.RoomId})
	if err != nil {
		if errors.Is(err, repository.UnknownRoomId) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("room does not exist"))
		} else {
			log.Logger.Errorf("could not get room: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not get room"))
		}
		return
	}

	var found bool
	for _, m := range room.Members {
		if m == body.Username {
			found = true
		}
	}

	if !found {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not part of room"))
		return
	}

	// check code
	var h *hub.RoomHub

	interfaceHub, ok := hubMap.Load(room.Id)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("room hub has not been loaded"))
		return
	} else {
		h = interfaceHub.(*hub.RoomHub)
	}

	code, err := h.SetHandshakeCode()
	if err != nil && errors.Is(err, hub.HandshakeInProgress) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("handshake is already in progress"))
		return
	}

	resp, err := json.Marshal(&startHandshakeResponse{
		Code: code,
	})
	if err != nil {
		log.Logger.Errorf("could not marshal handshake response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal handshake response"))
		return
	}

	w.Write(resp)
}
