package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"./shopify"
)

func validateDiscount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-control-allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`Method not found`))
		return
	}

	var body validateDiscountRequestBody
	bodyJSON, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(bodyJSON, &body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Couldn't parse request JSON:  %s", err)
		return
	}

	if len(body.Shop) < 1 || len(body.DiscountCode) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`The following parameters must not be empty: 'shop', 'discount_code'`))
		return
	}

	conf := getConfig()

	if _, exists := conf[body.Shop]; !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`Shop not authorized to use this API`))
		return
	}

	shop := shopify.New(body.Shop, conf[body.Shop])

	fmt.Printf("Request Data: %+v \n\n", body)

	discountCode, err := shop.LookupDiscountCode(body.DiscountCode)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Println(err)
		fmt.Fprintf(w, "false")
		return
	}

	priceRule, err := shop.GetPriceRule(discountCode.PriceRuleID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)
		return
	}

	fmt.Printf("Price Rule: %+v \n\n", priceRule)

	discountCodeValid := true

	if priceRule.EndDatePassed() {
		discountCodeValid = false
	}

	if priceRule.UsageLimit > 0 && discountCode.UsageCount >= priceRule.UsageLimit {
		discountCodeValid = false
	}

	//Check if priceRule entitled variant ids contains any variant ids sent in the request (if any)
	if len(body.VariantIDs) > 0 && len(priceRule.EntitledVariantIDs) > 0 {
		discountCodeValid = false
		for _, variantID := range body.VariantIDs {
			if !sliceContainsUint64(priceRule.EntitledVariantIDs, variantID) {
				continue
			}
			discountCodeValid = true
		}
	}

	//Check if priceRule entitled product ids contains any product ids sent in the request (if any)
	if len(body.ProductIDs) > 0 && len(priceRule.EntitledProductIDs) > 0 {
		discountCodeValid = false
		for _, productID := range body.ProductIDs {
			if !sliceContainsUint64(priceRule.EntitledProductIDs, productID) {
				continue
			}
			discountCodeValid = true
			break
		}
	}

	//Check if priceRule entitled collection ids contains any collection ids sent in the request (if any)
	if len(body.CollectionIDs) > 0 && len(priceRule.EntitledCollectionIDs) > 0 {
		discountCodeValid = false
		for _, collectionID := range body.CollectionIDs {
			if !sliceContainsUint64(priceRule.EntitledCollectionIDs, collectionID) {
				continue
			}
			discountCodeValid = true
			break
		}
	}

	respBody := strconv.FormatBool(discountCodeValid)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, respBody)
}
