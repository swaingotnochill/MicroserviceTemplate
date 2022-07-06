package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
)

// Product defines structure for an API product
type Product struct {
	ID           int     `json:"id"`
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description"`
	ProductPrice float32 `json:"price" validate:"gt=0"`
	SKU          string  `json:"sku" validate:"required,sku"`
	CreatedOn    string  `json:"_"`
	UpdatedOn    string  `json:"_"`
	DeletedOn    string  `json:"_"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// Adding validation to the Product structure
func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", SKUValidation)
	return validate.Struct(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	p.CreatedOn = time.Now().UTC().String()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	pos, err := GetProductById(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[pos] = p
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func GetProductById(id int) (int, error) {
	for i, p := range productList {
		if p.ID == id {
			return i, nil
		}
	}
	return -1, ErrProductNotFound
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var productList = []*Product{
	&Product{
		ID:           1,
		Name:         "Latte",
		Description:  "Frothy Milk Coffee",
		ProductPrice: 300,
		SKU:          "abs343",
		CreatedOn:    time.Now().UTC().String(),
		UpdatedOn:    time.Now().UTC().String(),
	},
	&Product{
		ID:           2,
		Name:         "Esspresso",
		Description:  "Short and strong coffee without milk",
		ProductPrice: 140,
		SKU:          "r2rd43",
		CreatedOn:    time.Now().UTC().String(),
		UpdatedOn:    time.Now().UTC().String(),
	},
}
