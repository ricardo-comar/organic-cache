package entity

type UserPricesEntity struct {
	UserId   string         `dynamodbav:"user_id" json:"user_id"`
	Products []ProductPrice `dynamodbav:"products" json:"products"`
}

type ProductPrice struct {
	ProductId     string  `dynamodbav:"product_id" json:"product_id"`
	ProductName   string  `dynamodbav:"product_name" json:"product_name"`
	OriginalValue float32 `dynamodbav:"original_value" json:"original_value"`
	Value         float32 `dynamodbav:"value" json:"value"`
	Discount      float32 `dynamodbav:"discount" json:"discount"`
}
