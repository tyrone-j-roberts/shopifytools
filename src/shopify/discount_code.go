package shopify

//DiscountCode model for parsing Shopify API response json
type DiscountCode struct {
	ID          uint64 `json:"id"`
	PriceRuleID uint64 `json:"price_rule_id"`
	Code        string `json:"code"`
	UsageCount  int    `json:"usage_count"`
}
