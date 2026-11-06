package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/api"
	"github.com/MaksimKrug/go-praktikum-final-project/pkg/config"
	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

const webDir = "./web"

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration: %w", err)
	}

	if err := db.Init(cfg.DBFile); err != nil {
		return fmt.Errorf("db init: %w", err)
	}
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init(cfg)

	log.Printf("server listening on port: %s", cfg.Port)
	if err := http.ListenAndServe(net.JoinHostPort("", cfg.Port), nil); err != nil {
		return fmt.Errorf("server: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
