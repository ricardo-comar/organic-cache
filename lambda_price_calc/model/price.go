package model

type UserPricesEntity struct {
	PriceId string         `dynamodbav:"id" json:"id"`
	UserId  string         `dynamodbav:"user_id" json:"user_id"`
	Prices  []ProductPrice `dynamodbav:"products" json:"products"`
}

type ProductPrice struct {
	ProductId     string  `json:"product_id"`
	OriginalValue float32 `json:"original_value"`
	Value         float32 `json:"value"`
	Discount      float32 `json:"discount"`
}
