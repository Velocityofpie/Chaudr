package main

import (
	"fmt"
	"net/http"
)

func createRoom() {
	fmt.Println("room created")
}

func leaveRoom() {
	fmt.Println("left room")
}

func main() {
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

		writer.Write([]byte("initiated handshake"))
	})

	// PUT /room
	// POST /room
	http.HandleFunc("/room", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPut && request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if request.Method == http.MethodPut {
			createRoom()
		} else if request.Method == http.MethodPost {
			leaveRoom()
		}

		writer.Write([]byte("hello world"))
	})

	http.HandleFunc("/greeting", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusBadRequest)
		}
		writer.Write([]byte("hello world"))
	})

	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		fmt.Println("server ran into an error: ", err)
	}
}
