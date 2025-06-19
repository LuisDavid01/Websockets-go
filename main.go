package main

import (
	//"context"
	//"github.com/LuisDavid01/Websockets-go/internal/manager"
	"log"
	"net/http"
	"time"

	"github.com/LuisDavid01/Websockets-go/internal/app"
	"github.com/LuisDavid01/Websockets-go/internal/routes"
)

func main() {
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	app.Logger.Println("the server started successfuly!:" + "https://localhost:8080/")

	log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
