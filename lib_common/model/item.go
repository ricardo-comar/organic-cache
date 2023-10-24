package model

type ProductItem struct {
	ProductId string  `dynamodbav:"id" json:"id"`
	Quantity  float32 `dynamodbav:"qtd" json:"qtd"`
}
