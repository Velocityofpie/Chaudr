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
var clientUI embed.FS

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	rand.Seed(time.Now().UnixNano())
	initDev()
	logger.Debug("running debug mode")

	serveMux := http.NewServeMux()

	flag.Parse()

	logger.Info("server starting...")
	err := http.ListenAndServe(*addr, addRoutes(serveMux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
