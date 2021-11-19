package routes

import (
	_ "embed"
	"encoding/json"
	hub2 "github.com/Velocityofpie/chaudr/hub"
	"github.com/Velocityofpie/chaudr/log"
	"github.com/pkg/errors"
	"math/rand"
	"net/http"
	"sync"
)

//go:embed test.html
var testHtml []byte

func AddRoutes(mux *http.ServeMux) http.Handler {
	// add dummy data
	dummyHub := hub2.NewHub()
	roomHubMap := new(sync.Map)
	roomHubMap.Store(uint(1234), dummyHub)
	go dummyHub.Run()

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
		hub2.BotHandler(roomHubMap, writer, request)
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
			log.Logger.Debug("creating room")
			hub := hub2.NewHub()
			roomId := rand.Uint32()
			roomHubMap.Store(roomId, &hub)
			go hub.Run()
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
		log.Logger.Infof(
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
