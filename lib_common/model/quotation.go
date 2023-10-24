package model

type ProductQuotation struct {
	ProductId     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Quantity      float32 `json:"qtd"`
	OriginalValue float32 `json:"original_value"`
	FinalValue    float32 `json:"final_value"`
	Discount      float32 `json:"discount"`
}
