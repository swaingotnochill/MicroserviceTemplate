package main

import (
	"context"
	""
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Dependency Injection
	//logger
	l := log.New(os.Stdout, "Product-API:", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/products", ph)

	// Basic Server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// Non-blocking go routine
	go func() {
		l.Println("Server Starting on PORT 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved Terminate, Graceful Shutdown. Signal Type: ", sig)

	// Graceful Shutdown
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
