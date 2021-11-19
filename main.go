package main

import (
	"embed"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//go:embed client/build/*
var client embed.FS

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

	log.Println("server starting...")
	err := http.ListenAndServe(*addr, serveMux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
