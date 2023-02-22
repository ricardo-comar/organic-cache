package model

type QuotationEntity struct {
	Id       string             `dynamodbav:"id" json:"id"`
	Products []ProductQuotation `dynamodbav:"products" json:"products"`
}

type ProductQuotation struct {
	ProductId     string  `json:"product_id"`
	Quantity      float32 `json:"qtd"`
	OriginalValue float32 `json:"original_value"`
	FinalValue    float32 `json:"final_value"`
	Discount      float32 `json:"discount"`
}
