package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/LuisDavid01/Websockets-go/internal/app"
	"github.com/LuisDavid01/Websockets-go/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "live-chat")
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	app.Logger.Printf("the server started successfuly on port: %d !!", port)

	log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
