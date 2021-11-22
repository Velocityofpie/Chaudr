package routes

import (
	"encoding/json"
	"github.com/Velocityofpie/chaudr/hub"
	"github.com/Velocityofpie/chaudr/log"
	"github.com/Velocityofpie/chaudr/repository"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"sync"
)

// PUT /room
// POST /room
func roomHandler(repo repository.RoomRepository, hubMap *sync.Map, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unexpected http method"))
		return
	}

	if r.Method == http.MethodPut {
		log.Logger.Debug("creating room")
		createRoom(repo, hubMap, w, r)
	} else if r.Method == http.MethodPost {
		log.Logger.Debug("leaving room")
		leaveRoom(repo, hubMap, w, r)
	}
}

type createRoomRequestBody struct {
	Username string `json:"username"`
}

type createRoomResponse struct {
	RoomId uint `json:"roomId"`
}

func createRoom(repo repository.RoomRepository, hubMap *sync.Map, w http.ResponseWriter, r *http.Request) {
	var body createRoomRequestBody
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

	if body.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body gave blank username"))
		return
	}

	createdRoom, err := repo.CreateRoom(repository.Room{Members: []string{body.Username}})
	if err != nil {
		log.Logger.Errorf("could not create room: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not create room"))
		return
	}

	h := hub.NewHub(createdRoom.Id)
	hubMap.Store(createdRoom.Id, h)
	go h.Run()

	if out, err := json.Marshal(createRoomResponse{RoomId: createdRoom.Id}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errors.Wrap(err, "failed to marshal response").Error()))
	} else {
		w.Write(out)
	}
}

type leaveRoomRequestBody struct {
	RoomId   uint   `json:"roomId"`
	Username string `json:"username"`
}

func leaveRoom(repo repository.RoomRepository, hubMap *sync.Map, w http.ResponseWriter, r *http.Request) {
	var body leaveRoomRequestBody
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

	room, err := repo.GetRoom(repository.Room{Id: body.RoomId})
	if err != nil {
		if errors.Is(err, repository.UnknownRoomId) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("request body gave invalid room id"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not get room details"))
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
		w.Write([]byte("unknown username for room"))
		return
	}

	_, err = repo.DeleteMemberFromRoom(room, body.Username)
	if err != nil {
		log.Logger.Errorf("could delete member %s from room %d: %v", body.Username, room.Id, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could delete member from room"))
		return
	}
}

type activeMembersRequestBody struct {
	RoomId   uint   `json:"roomId"`
	Username string `json:"username"`
}

type activeMembersResponse struct {
	Members []string `json:"members"`
}

// POST /room/members
func roomMembersHandler(repo repository.RoomRepository, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unexpected http method"))
		return
	}

	var body activeMembersRequestBody
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

	room, err := repo.GetRoom(repository.Room{Id: body.RoomId})
	if err != nil {
		if errors.Is(err, repository.UnknownRoomId) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("request body gave invalid room id"))
			return
		} else {
			log.Logger.Errorf("failed to get room: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errors.Wrap(err, "failed to get room").Error()))
		}
		return
	}

	resp := activeMembersResponse{Members: room.Members}
	if resp.Members == nil {
		resp.Members = make([]string, 0)
	}

	if out, err := json.Marshal(resp.Members); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errors.Wrap(err, "failed to marshal response").Error()))
	} else {
		w.Write(out)
	}
}
