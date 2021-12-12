package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MrDarkCoder/productapi/handlers"
)

// var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")
func main() {
	log.Println("WELCOME")
	// env.Parse()

	// create the handlers
	l := log.New(os.Stdout, "product_api", log.LstdFlags)
	ph := handlers.NewProducts(l)

	// create a new serve mux and register the handlers
	sermux := http.NewServeMux()
	sermux.Handle("/", ph)

	// manually creating a new server
	server := http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      sermux,            // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := server.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, err := context.WithTimeout(context.Background(), 30*time.Second)
	if err != nil {
		l.Print("went wrong quit the server")
	}
	server.Shutdown(ctx)
}
