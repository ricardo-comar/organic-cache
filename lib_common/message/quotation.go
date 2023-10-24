package message

import "github.com/ricardo-comar/organic-cache/lib_common/model"

type QuotationMessage struct {
	RequestId string                   `json:"request_id"`
	UserId    string                   `json:"user_id"`
	Products  []model.ProductQuotation `json:"products"`
}
