package model

type DiscountEntity struct {
	UserId    string            `dynamodbav:"user_id" json:"user_id"`
	Discounts []ProductDiscount `dynamodbav:"discounts" json:"discounts"`
}

type ProductDiscount struct {
	ProductId  string  `json:"product_id"`
	Percentage float32 `json:"percentage"`
}
