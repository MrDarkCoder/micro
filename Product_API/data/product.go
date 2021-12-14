package data

import "fmt"

// Product defines the structure for an API product
// swagger:model
type MyProduct struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"mysku"`
}

// Products Defines a slice of Product
type MyProducts []*MyProduct

// GetProducts : return all products from DB
func GetMyProducts() MyProducts {
	return myProductList
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func GetMyProductsByID(id int) (*MyProduct, error) {
	index := findIndexByProductID(id)
	if index == -1 || id == -1 {
		return nil, ErrProductNotFounded
	}
	return myProductList[index], nil
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func UpdateMyProduct(p *MyProduct) error {
	index := findIndexByProductID(p.ID)
	if index == -1 {
		return ErrProductNotFounded
	}

	// update the product in the DB
	myProductList[index] = p

	return nil
}

// AddProduct adds a new product to the database
func AddMyProduct(p *MyProduct) {
	// get the next id in sequence
	maxID := myProductList[len(myProductList)-1].ID
	p.ID = maxID + 1
	myProductList = append(myProductList, p)
}

// DeleteProduct deletes a product from the database
func DeleteMyProduct(id int) error {
	index := findIndexByProductID(id)
	if index == -1 {
		return ErrProductNotFound
	}

	myProductList = append(myProductList[:index], myProductList[index+1:]...)

	return nil
}

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFounded = fmt.Errorf("Product not found")

var myProductList = []*MyProduct{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
