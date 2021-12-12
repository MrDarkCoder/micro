package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	log.Println("WELCOME")

	http.HandleFunc("/", func(rw http.ResponseWriter, rq *http.Request) {
		log.Println("Hello World")
		data, err := ioutil.ReadAll(rq.Body)
		if err != nil {

			// rw.WriteHeader(http.StatusBadRequest)
			// rw.Write([]byte("went wrong"))

			// or use
			http.Error(rw, "went worng", http.StatusBadRequest)

			return
		}
		log.Printf("Data : %s\n", data)
		fmt.Fprintf(rw, "Hello %s\n", data)
	})

	http.HandleFunc("/goodbye", func(rw http.ResponseWriter, rq *http.Request) {
		log.Println("GoodBye World")
	})

	http.ListenAndServe(":9090", nil)

}
