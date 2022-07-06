package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"product-api/data"
	"product-api/handlers"

	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
)

// Global Variables
var bindAddresss = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {

	env.Parse()

	// Dependency Injection
	//logger
	l := log.New(os.Stdout, "Product-API:", log.LstdFlags)
	v := data.NewValidation()

	// Create the handlers
	ph := handlers.NewProducts(l, v)

	// Create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// Subrouter for handling GET requests
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.GetSingleProduct)

	// Subrouter for handling PUT requests
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddleWareValidationOfProduct) // This will run before the upper Handlefunc as this is a middleware

	// Subrouter for handling POST requests
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", ph.AddProducts)
	postRouter.Use(ph.MiddleWareValidationOfProduct) // This will run before the upper Handlefunc as this is a middleware

	// Basic Server
	s := &http.Server{
		Addr:         *bindAddresss,     // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		IdleTimeout:  120 * time.Second, // maximum time for connections using TCP Keep-Alive
		ReadTimeout:  1 * time.Second,   // maximum time for read request from the client
		WriteTimeout: 1 * time.Second,   // maximum time for write response to the client
	}

	// Non-blocking go routine for starting the server
	go func() {
		l.Println("Server Starting on PORT 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("[ERROR] starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// Trap the SIGTERM or INTERRUPT signal and gracefully shutdown the server
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block until we receive our signal
	sig := <-sigChan
	l.Println("Recieved Terminate, Graceful Shutdown. Signal Type: ", sig)

	// Graceful Shutdown the server, waiting max 30 seconds for current operations to complete
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
