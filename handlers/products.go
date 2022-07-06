// Package classification of Product API
//
// Documentation for Product API
//
//	Language: go
//  Schemes: http
//	BasePath: /products
// 	Version: 1.0.0
//
//	Consumes:
// 	- application/json
//	Produces:
//	- application/json
//swagger:meta
package handlers

import (
	"log"
	"net/http"
	"product-api/data"
	"strconv"

	"github.com/gorilla/mux"
)

// KeyProduct is a context key for product
type KeyProduct struct{}

// Products handler for getting and updating products
type Products struct {
	l *log.Logger
	v *data.Validation
}

// NewProducts returns a new products handler with given logger and validation
func NewProducts(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}


func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post products..")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("Product Data: %#v\n", prod)

	// save to datastore
	data.AddProduct(&prod)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	// Accessing variables passed int the URI
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to convert id to int", http.StatusBadRequest)
	}

	p.l.Println("Updating Products..")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Printf("Product Data: %#v\n", prod)

	// save to datastore
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Unable to update product", http.StatusInternalServerError)
		return
	}
}

