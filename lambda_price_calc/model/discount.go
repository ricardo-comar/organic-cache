package model

type DiscountEntity struct {
	UserId    string            `dynamodbav:"user_id" json:"user_id"`
	Discounts []ProductDiscount `dynamodbav:"discounts" json:"discounts"`
}

type ProductDiscount struct {
	Percentage float32 `dynamodbav:"percentage" json:"percentage"`
	ProductId  string  `dynamodbav:"product_id" json:"product_id"`
}
