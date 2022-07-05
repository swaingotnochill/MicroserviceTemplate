package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"product-api/data"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products..")

	// fetch data from datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusInternalServerError)
	}
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

type KeyProduct struct{}

func (p *Products) MiddleWareValidationOfProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{} // product structure from data.product
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}
		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx) // create a new context with the original context as upstream

		// call the next handler, which can be another middleware in the chain or the final handler.
		next.ServeHTTP(rw, r)
	})
}
