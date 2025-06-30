package main

import (
	"forum/internal/config"
	"forum/internal/router"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	app, err := router.NewApp(cfg)
	if err != nil {
		log.Fatal("failed to start app:", err)
	}

	log.Println("Server is running on http://localhost:" + cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, app); err != nil {
		log.Fatal(err)
	}
}
