package main

import (
	"log"
	"net/http"
)

type validateDiscountRequestBody struct {
	Shop          string   `json:"shop"`
	DiscountCode  string   `json:"discount_code"`
	VariantIDs    []uint64 `json:"variant_ids,omitempty"`
	ProductIDs    []uint64 `json:"product_ids,omitempty"`
	CollectionIDs []uint64 `json:"collection_ids,omitempty"`
}

func main() {
	http.HandleFunc("/validate-discount", validateDiscount)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
