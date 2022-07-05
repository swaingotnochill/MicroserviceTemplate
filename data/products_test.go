package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:         "Test",
		ProductPrice: 1.234,
		SKU:          "abc-dfg-Afb",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
