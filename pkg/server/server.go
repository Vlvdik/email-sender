package server

import (
	"log"
	"mail-Sender/config"
	"net/http"
)

type EmailServer struct {
	Server *http.Server
	Info   config.SendData
}

func pixelTracker(w http.ResponseWriter, r *http.Request) {
	log.Print("\n\nUser check email\n\nUser Email:", r.URL.Query().Get("email"))
	http.ServeFile(w, r, "server/pixelTracker.png")
}

func linkTracker(w http.ResponseWriter, r *http.Request) {
	log.Print("\n\nUser get to the link in letter\n\nUser Email:", r.URL.Query().Get("email"))
	http.Redirect(w, r, `Site with a "prize" for redirect`, http.StatusFound)
}

func (es *EmailServer) SetEmailServerData(sd config.SendData) {
	es.Info = sd
}

func (es *EmailServer) Start() error {
	es.Server = &http.Server{
		Addr: ":8082",
	}

	http.HandleFunc("/linkTracker", linkTracker)
	http.HandleFunc("/pixelTracker", pixelTracker)

	return es.Server.ListenAndServe()
}
