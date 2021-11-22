package routes

import (
	_ "embed"
	"github.com/Velocityofpie/chaudr/config"
	"github.com/Velocityofpie/chaudr/hub"
	"github.com/Velocityofpie/chaudr/log"
	"github.com/Velocityofpie/chaudr/repository"
	"net/http"
	"sync"

	"github.com/rs/cors"
)

//go:embed test.html
var testHtml []byte

func AddRoutes(mux *http.ServeMux) http.Handler {
	// create repositories
	repo, err := repository.NewSqliteRoomRepository()
	if err != nil {
		panic(err)
	}

	roomHubMap := new(sync.Map)

	// add dummy data
	if config.DebugMode {
		createdRoom, err := repo.CreateRoom(repository.Room{})
		if err != nil {
			log.Logger.Errorf("could not create debug room: %v", err)
		} else {
			dummyHub := hub.NewHub(createdRoom.Id)
			roomHubMap.Store(uint(1234), dummyHub)
			go dummyHub.Run()
			log.Logger.Debugf("created debug room: %d", createdRoom.Id)
		}
	}

	mux.HandleFunc("/room/connect", func(w http.ResponseWriter, r *http.Request) {
		joinRoomHandler(repo, roomHubMap, w, r)
	})

	// PUT /room/user
	mux.HandleFunc("/room/user", func(writer http.ResponseWriter, request *http.Request) {
		addUserToRoomHandler(repo, roomHubMap, writer, request)
	})

	// PUT /room/handshake
	mux.HandleFunc("/room/handshake", func(writer http.ResponseWriter, request *http.Request) {
		startHandshakeHandler(repo, roomHubMap, writer, request)
	})

	if config.DebugMode {
		// an api that adds a bot to an existing room which sends "hi" every ten seconds
		mux.HandleFunc("/room/hibot", func(writer http.ResponseWriter, request *http.Request) {
			hub.BotHandler(roomHubMap, writer, request)
		})

		mux.HandleFunc("/room/test", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write(testHtml)
		})
	}

	// PUT /room
	// POST /room
	mux.HandleFunc("/room", func(writer http.ResponseWriter, request *http.Request) {
		roomHandler(repo, roomHubMap, writer, request)
	})

	// POST /room/members
	mux.HandleFunc("/room/members", func(writer http.ResponseWriter, request *http.Request) {
		roomMembersHandler(repo, writer, request)
	})

	addHealthCheck(mux)

	if config.DebugMode {
		corsHandler := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedMethods:   []string{http.MethodPost, http.MethodPut, http.MethodGet},
			Debug:            true,
		}).Handler(mux)

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
			corsHandler.ServeHTTP(writer, request)
		})
	}

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
