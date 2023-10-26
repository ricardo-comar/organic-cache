package service

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/lib_common/model"
	"github.com/stretchr/testify/assert"
)

var ctx context.Context

var queryProductPriceCalled, queryRequestCalled, recalcMessageCalled, notifyQuotationCalled int
var pricesMsg *message.UserPricesMessage
var quotMsg *message.QuotationMessage
var pricesEntity *entity.UserPricesEntity
var quotationEntity *entity.QuotationEntity

func initVariables() {

	ctx = context.TODO()
	queryProductPriceCalled = 0
	queryRequestCalled = 0
	recalcMessageCalled = 0
	notifyQuotationCalled = 0

	pricesMsg = &message.UserPricesMessage{
		RequestId: uuid.New().String(),
		UserId:    "ABC",
	}

	pricesEntity = &entity.UserPricesEntity{
		UserId: pricesMsg.UserId,
		Products: []entity.ProductPrice{
			{ProductId: "A",
				ProductName:   "Mock Product A",
				OriginalValue: 100.0,
				Value:         90.0,
				Discount:      10.0},
		},
	}
	quotationEntity = &entity.QuotationEntity{
		RequestId: pricesMsg.RequestId,
		UserId:    pricesMsg.UserId,
		ProductList: []model.QuotationItem{
			{
				ProductId: pricesEntity.Products[0].ProductId,
				Quantity:  13.7,
			},
		},
	}

	quotMsg = &message.QuotationMessage{
		RequestId: pricesMsg.RequestId,
		UserId:    pricesMsg.UserId,
		Products: []model.ProductQuotation{
			{
				ProductId:     pricesEntity.Products[0].ProductId,
				Quantity:      quotationEntity.ProductList[0].Quantity,
				ProductName:   pricesEntity.Products[0].ProductName,
				OriginalValue: pricesEntity.Products[0].OriginalValue,
				Discount:      pricesEntity.Products[0].Discount,
				FinalValue:    pricesEntity.Products[0].Value * quotationEntity.ProductList[0].Quantity,
			},
		},
	}

}

func TestSuccessQuotation(t *testing.T) {
	initVariables()

	reqService := NewQuotationService(
		mockDynamoGateway{t: t, ctx: ctx, userId: pricesMsg.UserId, prices: pricesEntity, requestId: pricesMsg.RequestId, quot: quotationEntity},
		mockSQSGateway{t: t, ctx: ctx, quotation: quotMsg, quotationId: aws.String(uuid.New().String())},
	)

	reqErr := reqService.GenerateUserQuotation(ctx, pricesMsg)
	assert.Nil(t, reqErr, "Unexpected Error")

	assert.Equal(t, 1, queryProductPriceCalled, "Unexpected QueryProductPrice calls")
	assert.Equal(t, 1, queryRequestCalled, "Unexpected QueryRequest calls")
	assert.Equal(t, 0, recalcMessageCalled, "Unexpected RecalcMessage calls")
	assert.Equal(t, 1, notifyQuotationCalled, "Unexpected NotifyQuotation calls")

}
func TestSuccessRecalc(t *testing.T) {
	initVariables()

	reqService := NewQuotationService(
		mockDynamoGateway{t: t, ctx: ctx, userId: pricesMsg.UserId},
		mockSQSGateway{t: t, ctx: ctx, prices: pricesMsg, pricesId: aws.String(uuid.New().String())},
	)

	reqErr := reqService.GenerateUserQuotation(ctx, pricesMsg)
	assert.Nil(t, reqErr, "Unexpected Error")

	assert.Equal(t, 1, queryProductPriceCalled, "Unexpected QueryProductPrice calls")
	assert.Equal(t, 0, queryRequestCalled, "Unexpected QueryRequest calls")
	assert.Equal(t, 1, recalcMessageCalled, "Unexpected RecalcMessage calls")
	assert.Equal(t, 0, notifyQuotationCalled, "Unexpected NotifyQuotation calls")

}
func TestErrorQueryProductPrice(t *testing.T) {
	initVariables()

	mockError := errors.New("mock_error")

	reqService := NewQuotationService(
		mockDynamoGateway{t: t, ctx: ctx, userId: pricesMsg.UserId, errQueryProductPrice: mockError},
		mockSQSGateway{t: t, ctx: ctx},
	)

	reqErr := reqService.GenerateUserQuotation(ctx, pricesMsg)
	assert.Equal(t, mockError, reqErr, "Unexpected Error")

	assert.Equal(t, 1, queryProductPriceCalled, "Unexpected QueryProductPrice calls")
	assert.Equal(t, 0, queryRequestCalled, "Unexpected QueryRequest calls")
	assert.Equal(t, 0, recalcMessageCalled, "Unexpected RecalcMessage calls")
	assert.Equal(t, 0, notifyQuotationCalled, "Unexpected NotifyQuotation calls")

}
func TestErrorQueryRequest(t *testing.T) {
	initVariables()

	mockError := errors.New("mock_error")
	reqService := NewQuotationService(
		mockDynamoGateway{t: t, ctx: ctx, userId: pricesMsg.UserId, prices: pricesEntity, requestId: pricesMsg.RequestId, errQueryRequest: mockError},
		mockSQSGateway{t: t, ctx: ctx},
	)

	reqErr := reqService.GenerateUserQuotation(ctx, pricesMsg)
	assert.Equal(t, mockError, reqErr, "Unexpected Error")

	assert.Equal(t, 1, queryProductPriceCalled, "Unexpected QueryProductPrice calls")
	assert.Equal(t, 1, queryRequestCalled, "Unexpected QueryRequest calls")
	assert.Equal(t, 0, recalcMessageCalled, "Unexpected RecalcMessage calls")
	assert.Equal(t, 0, notifyQuotationCalled, "Unexpected NotifyQuotation calls")

}

func TestErrorRecalcMessage(t *testing.T) {
	initVariables()

	mockError := errors.New("mock_error")
	reqService := NewQuotationService(
		mockDynamoGateway{t: t, ctx: ctx, userId: pricesMsg.UserId},
		mockSQSGateway{t: t, ctx: ctx, prices: pricesMsg, errRecalcMessage: mockError},
	)

	reqErr := reqService.GenerateUserQuotation(ctx, pricesMsg)
	assert.Equal(t, mockError, reqErr, "Unexpected Error")

	assert.Equal(t, 1, queryProductPriceCalled, "Unexpected QueryProductPrice calls")
	assert.Equal(t, 0, queryRequestCalled, "Unexpected QueryRequest calls")
	assert.Equal(t, 1, recalcMessageCalled, "Unexpected RecalcMessage calls")
	assert.Equal(t, 0, notifyQuotationCalled, "Unexpected NotifyQuotation calls")

}

func TestErrorNotifyQuotation(t *testing.T) {
	initVariables()

	mockError := errors.New("mock_error")
	reqService := NewQuotationService(
		mockDynamoGateway{t: t, ctx: ctx, userId: pricesMsg.UserId, prices: pricesEntity, requestId: pricesMsg.RequestId, quot: quotationEntity},
		mockSQSGateway{t: t, ctx: ctx, quotation: quotMsg, errNotifyQuotation: mockError},
	)

	reqErr := reqService.GenerateUserQuotation(ctx, pricesMsg)
	assert.Equal(t, mockError, reqErr, "Unexpected Error")

	assert.Equal(t, 1, queryProductPriceCalled, "Unexpected QueryProductPrice calls")
	assert.Equal(t, 1, queryRequestCalled, "Unexpected QueryRequest calls")
	assert.Equal(t, 0, recalcMessageCalled, "Unexpected RecalcMessage calls")
	assert.Equal(t, 1, notifyQuotationCalled, "Unexpected NotifyQuotation calls")

}

type mockDynamoGateway struct {
	t         *testing.T
	ctx       context.Context
	userId    string
	requestId string
	prices    *entity.UserPricesEntity
	quot      *entity.QuotationEntity

	errQueryProductPrice error
	errQueryRequest      error
}

func (m mockDynamoGateway) QueryProductPrice(userId string) (*entity.UserPricesEntity, error) {
	queryProductPriceCalled++

	assert.Equal(m.t, m.ctx, ctx, "Unexpected ctx")
	assert.Equal(m.t, m.userId, userId, "Unexpected userId")

	return m.prices, m.errQueryProductPrice

}
func (m mockDynamoGateway) QueryRequest(requestId string) (*entity.QuotationEntity, error) {
	queryRequestCalled++

	assert.Equal(m.t, m.ctx, ctx, "Unexpected ctx")
	assert.Equal(m.t, m.requestId, requestId, "Unexpected userId")

	return m.quot, m.errQueryRequest
}

type mockSQSGateway struct {
	t                  *testing.T
	ctx                context.Context
	prices             *message.UserPricesMessage
	pricesId           *string
	errRecalcMessage   error
	quotation          *message.QuotationMessage
	quotationId        *string
	errNotifyQuotation error
}

func (m mockSQSGateway) RecalcMessage(ctx context.Context, msg *message.UserPricesMessage) (*string, error) {
	recalcMessageCalled++

	assert.Equal(m.t, m.ctx, ctx, "Unexpected ctx")
	if m.prices == nil {
		assert.Nil(m.t, msg, "msg is not nil")

	} else if diff := deep.Equal(m.prices, msg); diff != nil {
		m.t.Error(diff)
	}
	return m.pricesId, m.errRecalcMessage
}
func (m mockSQSGateway) NotifyQuotation(ctx context.Context, msg *message.QuotationMessage) (*string, error) {
	notifyQuotationCalled++

	assert.Equal(m.t, m.ctx, ctx, "Unexpected ctx")
	if m.quotation == nil {
		assert.Nil(m.t, msg, "msg is not nil")

	} else if diff := deep.Equal(m.quotation, msg); diff != nil {
		m.t.Error(diff)
	}

	return m.quotationId, m.errNotifyQuotation
}
