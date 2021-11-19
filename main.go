package main

import (
	"encoding/json"
	"flag"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// add dummy data
	dummyHub := newHub()
	roomHubMap := new(sync.Map)
	roomHubMap.Store(uint(1234), dummyHub)
	flag.Parse()

	http.HandleFunc("/", serveHome)

	http.HandleFunc("/room/connect", func(w http.ResponseWriter, r *http.Request) {
		log.Println("joining room")
		joinRoom(roomHubMap, w, r)
	})

	// PUT /room/user
	http.HandleFunc("/room/user", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPut {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.Write([]byte("added user to room"))
	})

	// PUT /room/handshake
	http.HandleFunc("/room/handshake", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPut {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.Write([]byte("needs to be implemented"))
	})

	// TODO: add an api that adds a bot to an existing room which sends "hi" every ten seconds

	// PUT /room
	// POST /room
	http.HandleFunc("/room", func(writer http.ResponseWriter, request *http.Request) {
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

	http.HandleFunc("/greeting", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusBadRequest)
		}
		writer.Write([]byte("hello world"))
	})

	log.Println("server starting...")
	err := http.ListenAndServe(*addr, http.DefaultServeMux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
