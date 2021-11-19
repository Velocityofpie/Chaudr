package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

func addRoutes(mux *http.ServeMux) http.Handler {
	// add dummy data
	dummyHub := newHub()
	roomHubMap := new(sync.Map)
	roomHubMap.Store(uint(1234), dummyHub)

	mux.HandleFunc("/room/connect", func(w http.ResponseWriter, r *http.Request) {
		log.Println("joining room")
		joinRoom(roomHubMap, w, r)
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

	// TODO: add an api that adds a bot to an existing room which sends "hi" every ten seconds

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
