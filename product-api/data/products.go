package data

import (
	"encoding/json"
	"io"
	"time"
)

// Product defines structure for an API product
type Product struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ProductPrice float32 `json:"price"`
	SKU          string  `json:"sku"`
	CreatedOn    string  `json:"_"`
	UpdatedOn    string  `json:"_"`
	DeletedOn    string  `json:"_"`
}

func (p *Products) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return productList
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
		ID:           1,
		Name:         "Esspresso",
		Description:  "Short and strong coffee without milk",
		ProductPrice: 140,
		SKU:          "r2rd43",
		CreatedOn:    time.Now().UTC().String(),
		UpdatedOn:    time.Now().UTC().String(),
	},
}
