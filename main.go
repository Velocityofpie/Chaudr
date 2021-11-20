package main

import (
	"embed"
	"flag"
	"github.com/Velocityofpie/chaudr/log"
	"github.com/Velocityofpie/chaudr/routes"
	"io/fs"
	"math/rand"
	"net/http"
	"time"
)

//go:embed client/build/*
var clientUI embed.FS

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Logger.Debug("running debug mode")

	serveMux := http.NewServeMux()
	// add javascript client
	f, err := fs.Sub(clientUI, "client/build")
	if err != nil {
		panic(err)
	}
	serveMux.Handle("/", http.FileServer(http.FS(f)))
	flag.Parse()

	log.Logger.Info("server starting...")
	err = http.ListenAndServe(*addr, routes.AddRoutes(serveMux))
	if err != nil {
		log.Logger.Fatal("ListenAndServe: ", err)
	}
}
