package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	serverName := os.Getenv("SERVER_NAME")
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	if serverName == "" {
		serverName = "localhost"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("[%s]Request received", serverName)))
	})

	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
