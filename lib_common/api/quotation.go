package api

import "github.com/ricardo-comar/organic-cache/lib_common/model"

type QuotationRequest struct {
	RequestId string                `json:"request_id"`
	UserId    string                `json:"user_id"`
	Products  []model.QuotationItem `json:"products"`
}

type QuotationResponse struct {
	RequestId string                   `json:"request_id"`
	UserId    string                   `json:"user_id"`
	Products  []model.ProductQuotation `json:"products"`
}
