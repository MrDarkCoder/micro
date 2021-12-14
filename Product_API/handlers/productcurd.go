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

// KeyProduct is a key used for the Product object in the context
type MyKeyProduct struct{}

// Products handler for getting and updating products
type MyProducts struct {
	l *log.Logger
	v *data.Validation
}

// NewProducts returns a new products handler with the given logger
func NewMyProducts(l *log.Logger, v *data.Validation) *MyProducts {
	return &MyProducts{l, v}
}

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("invalid Path, path should be /products/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getMyProductID(rq *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(rq)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all current products
func (p *MyProducts) ListAll(rw http.ResponseWriter, rq *http.Request) {
	p.l.Println("[DEBUG] Get All Products")

	prods := data.GetMyProducts()

	err := data.TooJSON(prods, rw)

	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}

}

// swagger:route GET /products/{id} products listSingle
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *MyProducts) ListSingle(rw http.ResponseWriter, rq *http.Request) {
	id := getMyProductID(rq)
	p.l.Println("[DEBUG] get record id", id)

	prod, err := data.GetMyProductsByID(id)

	switch err {
	case nil: // passed
	case data.ErrProductNotFounded:
		p.l.Println("[ERROR] fetching product", err)
		rw.WriteHeader(http.StatusNotFound)
		data.TooJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.TooJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.TooJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}

}

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new products
func (p *MyProducts) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(*data.MyProduct)

	p.l.Printf("[DEBUG] Inserting product: %#v\n", prod)
	data.AddMyProduct(prod)
}

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
func (p *MyProducts) Update(rw http.ResponseWriter, rq *http.Request) {

	// fetch the product from the context
	prod := rq.Context().Value(KeyProduct{}).(*data.MyProduct)
	p.l.Println("[DEBUG] updating record id", prod.ID)

	err := data.UpdateMyProduct(prod)
	if err == data.ErrProductNotFounded {
		p.l.Println("[ERROR] product not found", err)

		rw.WriteHeader(http.StatusNotFound)
		data.TooJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}

// swagger:route DELETE /products/{id} products deleteProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (p *MyProducts) Delete(rw http.ResponseWriter, rq *http.Request) {
	id := getMyProductID(rq)

	p.l.Println("[DEBUG] deleting record id", id)

	err := data.DeleteMyProduct(id)
	if err == data.ErrProductNotFounded {
		p.l.Println("[ERROR] deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.TooJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.TooJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// MiddlewareValidateMyProduct validates the product in the request and calls next if ok
func (p *MyProducts) MiddlewareValidateMyProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		prod := &data.MyProduct{}

		err := data.FroomJSON(prod, rq.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.TooJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		// validate the product
		errs := p.v.Validate(prod)
		if errs != nil || len(errs) != 0 {
			p.l.Println("[ERROR] validating product", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.TooJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add the product to the context
		ctx := context.WithValue(rq.Context(), KeyProduct{}, prod)
		rq = rq.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, rq)
	})
}
