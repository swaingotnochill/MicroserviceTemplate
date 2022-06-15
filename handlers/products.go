package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProducts(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		// // expect id in the URI
		// r.URL.Query().Get("id")
		p.l.Println("PUT", r.URL.Path)

		regex := regexp.MustCompile(`/([0-9]+)`)
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)
		p.l.Println("g", g)
		p.l.Println("g[0]", g[0])

		if len(g) != 1 {
			p.l.Println("Invalid URI, more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI, more than two capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.l.Println("Got id :", id)

		p.updateProducts(id, rw, r)
		return
	}

	// Catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products..")

	// fetch data from datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusInternalServerError)
	}
}

func (p *Products) addProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post products..")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to Unmarshal JSON", http.StatusBadRequest)
	}

	p.l.Printf("Product Data: %#v\n", prod)

	// save to datastore
	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Updating Products..")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to Unmarshal JSON", http.StatusBadRequest)
	}

	p.l.Printf("Product Data: %#v\n", prod)

	// save to datastore
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Unable to update product", http.StatusInternalServerError)
		return
	}
}
