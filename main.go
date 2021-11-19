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

	serveMux := http.NewServeMux()

	flag.Parse()

	log.Println("server starting...")
	err := http.ListenAndServe(*addr, addRoutes(serveMux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
