package model

type QuotationEntity struct {
	Id       string             `dynamodbav:"id" json:"id"`
	Products []ProductQuotation `dynamodbav:"products" json:"products"`
	TTL      string             `dynamodbav:"ttl" json:"ttl"`
}

type ProductQuotation struct {
	ProductId     string  `dynamodbav:"product_id" json:"product_id"`
	Quantity      float32 `dynamodbav:"qtd" json:"qtd"`
	OriginalValue float32 `dynamodbav:"original_value" json:"original_value"`
	FinalValue    float32 `dynamodbav:"final_" json:"final_value"`
	Discount      float32 `dynamodbav:"discount" json:"discount"`
}
