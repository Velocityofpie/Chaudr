package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func createRoom() {
	fmt.Println("room created")
}

func leaveRoom() {
	fmt.Println("left room")
}

//func main() {
//	// PUT /room/user
//	http.HandleFunc("/room/user", func(writer http.ResponseWriter, request *http.Request) {
//		if request.Method != http.MethodPut {
//			writer.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		writer.Write([]byte("added user to room"))
//	})
//
//	// PUT /room/handshake
//	http.HandleFunc("/room/handshake", func(writer http.ResponseWriter, request *http.Request) {
//		if request.Method != http.MethodPut {
//			writer.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		writer.Write([]byte("initiated handshake"))
//	})
//
//	// PUT /room
//	// POST /room
//	http.HandleFunc("/room", func(writer http.ResponseWriter, request *http.Request) {
//		if request.Method != http.MethodPut && request.Method != http.MethodPost {
//			writer.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		if request.Method == http.MethodPut {
//			createRoom()
//		} else if request.Method == http.MethodPost {
//			leaveRoom()
//		}
//
//		writer.Write([]byte("hello world"))
//	})
//
//	http.HandleFunc("/greeting", func(writer http.ResponseWriter, request *http.Request) {
//		if request.Method != http.MethodGet {
//			writer.WriteHeader(http.StatusBadRequest)
//		}
//		writer.Write([]byte("hello world"))
//	})
//
//	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
//		fmt.Println("server ran into an error: ", err)
//	}
//}

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
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
