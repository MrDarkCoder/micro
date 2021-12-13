package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	// handle the request for a CREATE A PRODUCT
	if rq.Method == http.MethodPost {
		p.addProduct(rw, rq)
		return
	}

	// handle the request for update
	if rq.Method == http.MethodPut{
		p.l.Println("PUT", rq.URL.Path)
		// expect the id in the URI
		r := `/([0-9]+)`
		reg := regexp.MustCompile(r)
		g := reg.FindAllStringSubmatch(rq.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(rw, "Invaild URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invaild URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert to numer", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.updateProducts(id, rw, rq)
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

func (p *Products) addProduct(rw http.ResponseWriter, rq *http.Request) {
	p.l.Println("Handle POST Product")

	prod := &data.Product{}

	err := prod.FromJSON(rq.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p Products) updateProducts(id int, rw http.ResponseWriter, rq *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}

	err := prod.FromJSON(rq.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}