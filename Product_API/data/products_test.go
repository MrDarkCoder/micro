package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		ID:          0,
		Name:        "nave",
		Description: "",
		Price:       1.00,
		SKU:         "abs-abc-def",
		CreatedOn:   "",
		UpdatedOn:   "",
		DeletedOn:   "",
	}
	err := p.Validate()
	if err != nil{
		t.Fatal(err)
	}
}