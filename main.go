package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

const (
	webDir        = "./web"
	portDefault   = "7540"
	dbFileDefault = "scheduler.db"
)

func getPort() (string, error) {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = portDefault
	}

	portParsed, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("port %q is not a number: %w", port, err)
	}

	if portParsed > 65535 || portParsed < 1 {
		return "", fmt.Errorf("port is outside [1, 65535]: %s", port)
	}

	return port, nil
}

func getDBPath() string {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = dbFileDefault
	}

	return dbPath
}

func main() {

	// create db if needed
	dbPath := getDBPath()
	err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("db Init error: %v", err)
	}
	defer db.Close()

	// handle web-server
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port, err := getPort()
	if err != nil {
		log.Fatalf("getPort error: %v", err)
	}

	log.Printf("Server Listening on port: %s", port)
	if err := http.ListenAndServe(net.JoinHostPort("", port), nil); err != nil {
		log.Fatalf("server not started: %v", err)
	}

}
