package handlers

import (
	"net/http"
	"product-api/data"
	"strconv"
)

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

func (p *Products) GetSingleProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET single product..")
	idFromURL, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		http.Error(rw, "Unable to convert id to int", http.StatusBadRequest)
		return
	}

	prod, err := data.GetProductById(idFromURL)

	if err != nil {
		http.Error(rw, "[Error] Product not found", http.StatusBadRequest)
		return
	}
	p.l.Println("[DEBUG] get record id", idFromURL)

	err = data.ToJSON(prod, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
	}

}
