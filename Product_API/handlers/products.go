package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MrDarkCoder/productapi/data"
	"github.com/gorilla/mux"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, rq *http.Request) {
	// fetch the products from the datastore
	lp := data.GetProducts()

	// d, err := json.Marshal(lp) instead if this, use encoders(lit faster)

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, rq *http.Request) {
	p.l.Println("Handle POST Product")
	// using our middleware
	prod := rq.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)

	/*
		prod := &data.Product{}
		err := prod.FromJSON(rq.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		}
		data.AddProduct(prod)
	*/
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, rq *http.Request) {
	p.l.Println("Handle PUT Product")

	// getting id using mux
	params := mux.Vars(rq)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(rw, "Unable to COnvert to string", http.StatusBadRequest)
		return
	}

	/*
		prod := &data.Product{}

		err := prod.FromJSON(rq.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		}
		err = data.UpdateProduct(id, prod)
	*/

	p.l.Println("Handle PUT Product", id)
	prod := rq.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

// Delete Handler
// func () DeleteProduct(){}

// middleware for product validation
type KeyProduct struct{}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		// take pur product
		prod := data.Product{}

		// Deserialize
		err := prod.FromJSON(rq.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the product to the context
		ctx := context.WithValue(rq.Context(), KeyProduct{}, prod)
		rq = rq.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, rq)
	})
}

/*
ServeHTTP is the main entry point for the handler and staisfies the http.Handler
interface
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
	rw.WriteHeader(http.StatusMethodNotAllowed)}
*/
