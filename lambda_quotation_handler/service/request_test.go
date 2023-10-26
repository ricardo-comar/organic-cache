package service_test

import (
	"testing"

	"github.com/ricardo-comar/organic-cache/quotation_handler/service"
	"github.com/stretchr/testify/assert"
)

func TestSuccessRequest(t *testing.T) {
	initVariables()

	reqService := service.NewRequestService(mockDynamoGateway{t: t, req: quotEnt}, mockSNSGateway{t: t, ctx: ctx, msg: snsMsg, messageId: &msgId})

	reqErr := reqService.RequestQuotation(ctx, quotReq)
	assert.Nil(t, reqErr, "Unexpected Error")
}
func TestErrorSaveQuotationRequest(t *testing.T) {
	initVariables()

	reqService := service.NewRequestService(mockDynamoGateway{t: t, req: quotEnt, err: err}, mockSNSGateway{t: t, ctx: ctx, msg: snsMsg, messageId: &msgId})

	reqErr := reqService.RequestQuotation(ctx, quotReq)
	assert.Equal(t, err, reqErr, "Unexpected Error")
}

func TestErrorNotifyQuotation(t *testing.T) {
	initVariables()

	reqService := service.NewRequestService(mockDynamoGateway{t: t, req: quotEnt}, mockSNSGateway{t: t, ctx: ctx, msg: snsMsg, err: err})

	reqErr := reqService.RequestQuotation(ctx, quotReq)
	assert.Equal(t, err, reqErr, "Unexpected Error")
}
