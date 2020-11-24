package main

import (
	"flag"
	"log"
	"net/http"

	"go_console_chat/controller"
	"go_console_chat/model"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	hub := model.NewHub()
	go hub.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controller.WSHandler(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
