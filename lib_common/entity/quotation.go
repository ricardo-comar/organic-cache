package entity

import "github.com/ricardo-comar/organic-cache/lib_common/model"

type QuotationEntity struct {
	RequestId   string                `dynamodbav:"request_id" json:"request_id"`
	UserId      string                `dynamodbav:"user_id" json:"user_id"`
	ProductList []model.QuotationItem `dynamodbav:"products" json:"products"`
	TTL         string                `dynamodbav:"ttl" json:"ttl"`
}
