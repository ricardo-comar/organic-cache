package entity

type ProductEntity struct {
	ProductId string  `dynamodbav:"id" json:"id"`
	Name      string  `dynamodbav:"name" json:"name"`
	Value     float32 `dynamodbav:"value" json:"value"`
}
