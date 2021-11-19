package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/Velocityofpie/chaudr/repository"
	"github.com/pkg/errors"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

//go:embed test.html
var testHtml []byte

// botHandler handles websocket requests from the peer.
func botHandler(hubMap *sync.Map, w http.ResponseWriter, r *http.Request) {
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

	for client := range h.clients {
		if client.member.Username == username {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("hibot username is already taken"))
			return
		}
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
	hub, ok := hubMap.Load(uint(id))
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unknown room id %d", id)))
		return
	}

	h := hub.(*RoomHub)

	for client := range h.clients {
		if client.member.Username == username {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("member is already connected: %s", username)))
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &ConnectedMember{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
		member: repository.Member{
			RoomID:   uint(id),
			Username: username,
		},
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func addRoutes(mux *http.ServeMux) http.Handler {
	// add dummy data
	dummyHub := newHub()
	roomHubMap := new(sync.Map)
	roomHubMap.Store(uint(1234), dummyHub)
	go dummyHub.run()

	mux.HandleFunc("/room/connect", func(w http.ResponseWriter, r *http.Request) {
		joinRoomHandler(roomHubMap, w, r)
	})

	// PUT /room/user
	mux.HandleFunc("/room/user", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPut {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.Write([]byte("added user to room"))
	})

	// PUT /room/handshake
	mux.HandleFunc("/room/handshake", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPut {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.Write([]byte("needs to be implemented"))
	})

	// an api that adds a bot to an existing room which sends "hi" every ten seconds
	mux.HandleFunc("/room/hibot", func(writer http.ResponseWriter, request *http.Request) {
		botHandler(roomHubMap, writer, request)
	})

	mux.HandleFunc("/room/test", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write(testHtml)
	})

	// PUT /room
	// POST /room
	mux.HandleFunc("/room", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPut && request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if request.Method == http.MethodPut {
			// TODO: create room logic
			log.Println("creating room")
			hub := newHub()
			roomId := rand.Uint32()
			roomHubMap.Store(roomId, &hub)
			go hub.run()
			type createRoomResponse struct {
				RoomId uint `json:"roomId"`
			}

			if out, err := json.Marshal(createRoomResponse{RoomId: uint(roomId)}); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(errors.Wrap(err, "failed to marshal response").Error()))
			} else {
				writer.Write(out)
			}
			return

		} else if request.Method == http.MethodPost {
			// TODO: leave room logic
		}

	})

	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusBadRequest)
		}
		writer.Write([]byte("healthy"))
	})

	// add javascript client
	f, err := fs.Sub(clientUI, "client/build")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(f)))

	// add logging middleware
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		scheme := "http"
		if request.TLS != nil {
			scheme = "https"
		}
		logger.Infof(
			"%s://%s%s %s from %s",
			scheme,
			request.Host,
			request.RequestURI,
			request.Proto,
			request.RemoteAddr,
		)
		mux.ServeHTTP(writer, request)
	})
}
