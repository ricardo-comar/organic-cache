package model

type QuotationRequest struct {
	RequestId    string        `dynamodbav:"request_id" json:"request_id"`
	ConnectionId string        `dynamodbav:"connection_id" json:"connection_id"`
	UserId       string        `dynamodbav:"user_id" json:"user_id"`
	ProductList  []ProductItem `dynamodbav:"products" json:"products"`
	TTL          string        `dynamodbav:"ttl" json:"ttl"`
}

type ProductItem struct {
	ProductId string  `dynamodbav:"id" json:"id"`
	Quantity  float32 `dynamodbav:"qtd" json:"qtd"`
}
