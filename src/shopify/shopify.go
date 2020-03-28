package shopify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//New initialised and returns new instance of Shopify
func New(shop, password string) Shopify {
	return Shopify{BaseURL: fmt.Sprintf("https://%s", shop), Password: password}
}

//Shopify Object, to use with API requests
type Shopify struct {
	BaseURL  string
	Password string
}

func handleShopifyCallLimit(headers map[string][]string) {
	if headerVals, exists := headers["Http_x_shopify_shop_api_call_limit"]; exists && len(headerVals) > 0 {
		limitParts := strings.Split(headerVals[0], "/")
		limitCurrent, _ := strconv.Atoi(limitParts[0])
		limitMax, _ := strconv.Atoi(limitParts[1])

		fmt.Printf("Shopify Call Limit: %d / %d \n", limitCurrent, limitMax)

		if limitCurrent >= limitMax-1 {
			fmt.Println("Approaching shopify call limit. Sleeping...")
			time.Sleep(20 * time.Second)
		}
	}
}

func (shop Shopify) makeRequest(method, reqURL string, jsonRequest []byte) (*http.Response, error) {
	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	var reqBody io.Reader

	if jsonRequest != nil {
		reqBody = bytes.NewReader(jsonRequest)
	}

	bytes.NewReader(jsonRequest)

	req, err := http.NewRequest(method, reqURL, reqBody)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Shopify-Access-Token", shop.Password)

	return client.Do(req)
}

//LookupDiscountCode uses the Shopify DiscountCode API to get Discount code data by the code string
func (shop Shopify) LookupDiscountCode(code string) (DiscountCode, error) {
	reqURL := shop.BaseURL + "/admin/api/2020-01/discount_codes/lookup.json?code=" + code
	resp, err := shop.makeRequest(http.MethodGet, reqURL, nil)

	fmt.Println(reqURL)

	var discountCode DiscountCode

	if err != nil {
		return discountCode, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)

	if resp.StatusCode != http.StatusSeeOther {
		return discountCode, errors.New(string(body))
	}

	reqURL = fmt.Sprintf("%s.json", resp.Header["Location"][0])
	resp, err = shop.makeRequest(http.MethodGet, reqURL, nil)

	if err != nil {
		return discountCode, err
	}

	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)

	response := struct {
		DiscountCode DiscountCode `json:"discount_code"`
	}{DiscountCode: discountCode}

	json.Unmarshal(body, &response)

	return response.DiscountCode, nil
}

//GetPriceRule uses the Shopify PriceRule API endpoint to get PriceRule data by the id
func (shop Shopify) GetPriceRule(price_rule_id uint64) (PriceRule, error) {
	reqURL := fmt.Sprintf("%s/admin/api/2020-01/price_rules/%d.json", shop.BaseURL, price_rule_id)
	resp, err := shop.makeRequest(http.MethodGet, reqURL, nil)

	var priceRule PriceRule

	if err != nil {
		return priceRule, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return priceRule, errors.New(string(body))
	}

	response := struct {
		PriceRule PriceRule `json:"price_rule"`
	}{PriceRule: priceRule}

	json.Unmarshal(body, &response)

	return response.PriceRule, nil
}
