package main

import (
	"fmt"
	"log"
	"ws_implementation/db"
	"ws_implementation/internal/user"
	"ws_implementation/internal/ws"
	"ws_implementation/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize db connection: %v", err)
	}

	repo := user.NewRepository(dbConn.GetDB())
	svc := user.NewService(repo)
	userHandler := user.NewHandler(svc)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	router.InitRouter(userHandler, wsHandler)
	router.Start(":8080")
	fmt.Println("it is listening")
}
