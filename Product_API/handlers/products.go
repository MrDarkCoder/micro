package handlers

import (
	"log"
	"net/http"

	"github.com/MrDarkCoder/productapi/data"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}
// ServeHTTP is the main entry point for the handler and staisfies the http.Handler
// interface
func (p *Products) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
		// handle the request for a list of products
		if rq.Method == http.MethodGet {
			p.getProducts(rw, rq)
			return
		}
	
		// catch all
		// if no method is satisfied return an error
		rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, rq *http.Request) {
	// fetch the products from the datastore
	lp := data.GetProducts()
	
	// d, err := json.Marshal(lp) instead if this, use encoders(lit faster)

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
