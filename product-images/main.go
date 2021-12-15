package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MrDarkCoder/productimages/files"
	"github.com/MrDarkCoder/productimages/handlers"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("WELCOME")
	l := log.New(os.Stdout, "product_api", log.LstdFlags)

	// create the storage class, use local storage
	// max filesize 5MB
	stor, err := files.NewLocal("./imagestore", 1024*1000*5)
	if err != nil {
		l.Println("Unable to create storage", "error", err)
		os.Exit(1)
	}

	fh := handlers.NewFiles(stor, l)

	sermux := mux.NewRouter()

	// filename regex: {filename:[a-zA-Z]+\\.[a-z]{3}}
	// problem with FileServer is that it is dumb
	ph := sermux.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.ServeHTTP)

	// get files
	gh := sermux.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir("./imagestore"))),
	)

	// create a new server
	server := http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      sermux,            // set the default handler
		ErrorLog:     l,                 // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server", "bind_address localhost:9090")

		err := server.ListenAndServe()
		if err != nil {
			l.Fatal("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Println("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

// curl -vv localhost:9090/1/go.mod -X PUT --data-binary @test.png
