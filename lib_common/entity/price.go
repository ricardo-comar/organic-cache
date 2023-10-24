package entity

type UserPricesEntity struct {
	UserId   string         `dynamodbav:"user_id" json:"user_id"`
	Products []ProductPrice `dynamodbav:"products" json:"products"`
}

type ProductPrice struct {
	ProductId     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	OriginalValue float32 `json:"original_value"`
	Value         float32 `json:"value"`
	Discount      float32 `json:"discount"`
}
