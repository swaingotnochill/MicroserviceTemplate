package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"product-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Dependency Injection
	//logger
	l := log.New(os.Stdout, "Product-API:", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()

	// Subrouter for handling GET requests
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)

	// Subrouter for handling PUT requests
	putRouter := sm.Methods(http.MethodPut).Subrouter()

	// defining id variable in the URI which can be accessed by handler using mux.vars which is a map of variables passed in the URI
	putRouter.HandleFunc("/products/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddleWareValidationOfProduct) // This will run before the upper Handlefunc as this is a middleware

	// Subrouter for handling POST requests
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", ph.AddProducts)
	postRouter.Use(ph.MiddleWareValidationOfProduct) // This will run before the upper Handlefunc as this is a middleware

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
