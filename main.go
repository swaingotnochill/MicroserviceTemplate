package main

import (
	"example/hello/handlers"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Dependency Injection
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	hh := handlers.NewHello(l)
	gh := handlers.NewBye(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/bye", gh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	s.ListenAndServe()
}
