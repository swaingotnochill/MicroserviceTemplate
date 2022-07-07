package handlers

import (
	"context"
	"fmt"
	"net/http"
	"product-api/data"
)

// MiddlewareVlidationOfProduct validates the product in the request and calls next if ok
func (p *Products) MiddleWareValidationOfProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{} // product structure from data.product

		err := data.FromJSON(prod, r.Body)
		if err != nil {
			http.Error(rw, "[Error] reading product", http.StatusBadRequest)
			return
		}

		// validate the product before adding it to the context
		err = prod.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("[Error] validating product: %s", err), http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx) // create a new context with the original context as upstream

		// call the next handler, which can be another middleware in the chain or the final handler.
		next.ServeHTTP(rw, r)
	})
}
