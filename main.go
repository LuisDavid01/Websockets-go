package main

import (
	"context"
	"log"
	"net/http"
)

func main() {

	ctx := context.Background()

	manager := NewManager(ctx)

	http.Handle(("/"), http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.ServeWS)
	http.HandleFunc("/login", manager.loginHandler)

	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
}
