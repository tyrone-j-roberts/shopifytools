package shopify

import "time"

//PriceRule model for parsing Shopify API response json
type PriceRule struct {
	ID                        uint64         `json:"id"`
	ValueType                 string         `json:"value_type"`
	Value                     int            `json:"value"`
	CustomerSelection         string         `json:"customer_selection"`
	TargetType                string         `json:"target_type"`
	TargetSelection           string         `json:"target_selection"`
	StartsAt                  string         `json:"starts_at"`
	EndsAt                    string         `json:"ends_at"`
	UsageLimit                int            `json:"usage_limit"`
	EntitledProductIDs        []uint64       `json:"entitled_product_ids"`
	EntitledVariantIDs        []uint64       `json:"entitled_variant_ids"`
	EntitledCollectionIDs     []uint64       `json:"entitled_collection_ids"`
	PrerequisiteQuantityRange map[string]int `json:"prerequisite_quantity_range"`
}

//EndDatePassed returns true is current date is after PriceRule end date, otherwise returns false
func (rule PriceRule) EndDatePassed() bool {
	now := time.Now()
	if len(rule.EndsAt) < 1 {
		return false
	}
	endDate, _ := time.Parse(time.RFC3339, rule.EndsAt)
	return now.After(endDate)
}
