package model

type ProductEntity struct {
	PriceId string  `dynamodbav:"id" json:"id"`
	Value   float32 `dynamodbav:"value" json:"value"`
}
