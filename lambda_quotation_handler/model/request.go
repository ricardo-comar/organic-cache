package model

//EmployeeRegistries slice
type QuotationRequest struct {
	UserId      string        `json:"id"`
	ProductList []ProductItem `json:"products"`
}

type ProductItem struct {
	ProductId string  `json:"id"`
	Quantity  float32 `json:"qtd"`
}
