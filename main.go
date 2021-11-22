package main

import (
	"context"
	"embed"
	"flag"
	"github.com/Velocityofpie/chaudr/config"
	"github.com/Velocityofpie/chaudr/log"
	"github.com/Velocityofpie/chaudr/routes"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//go:embed client/src/client/dist/*
var clientUI embed.FS

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Logger.Debug("running debug mode")

	serveMux := http.NewServeMux()
	// add javascript client
	f, err := fs.Sub(clientUI, "client/src/client/dist")
	if err != nil {
		panic(err)
	}
	serveMux.Handle("/", http.FileServer(http.FS(f)))
	flag.Parse()

	srv := &http.Server{Addr: *addr, Handler: routes.AddRoutes(serveMux)}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if config.DebugMode {
			log.Logger.Debug("delete database")
			os.Remove("gorm.db")
		}
		srv.Shutdown(context.Background())
	}()

	log.Logger.Info("server starting...")
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Logger.Fatal("ListenAndServe: ", err)
	}
}
