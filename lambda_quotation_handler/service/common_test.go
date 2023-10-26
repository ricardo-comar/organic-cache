package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"github.com/ricardo-comar/organic-cache/lib_common/api"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/lib_common/model"
	"github.com/stretchr/testify/assert"
)

type mockDynamoGateway struct {
	t   *testing.T
	req entity.QuotationEntity
	err error
}

func (g mockDynamoGateway) SaveQuotationRequest(req entity.QuotationEntity) error {
	assert.Equal(g.t, g.req.ProductList, req.ProductList, "Empty ProducList")
	assert.Equal(g.t, g.req.RequestId, req.RequestId, "Empty RequestId")
	assert.Equal(g.t, g.req.UserId, req.UserId, "Empty UserId")
	assert.NotEmpty(g.t, req.TTL, "Empty TTL")

	return g.err
}

type mockSNSGateway struct {
	t         *testing.T
	ctx       context.Context
	msg       message.UserPricesMessage
	messageId *string
	err       error
}

func (g mockSNSGateway) NotifyQuotation(ctx context.Context, msg message.UserPricesMessage) (*string, error) {
	assert.Equal(g.t, g.ctx, ctx, "Unexpected context")
	assert.Equal(g.t, g.msg.RequestId, msg.RequestId, "Invalid RequestId")
	assert.Equal(g.t, g.msg.UserId, msg.UserId, "Invalid UserId")

	return g.messageId, g.err
}

type mockSQSGateway struct {
	t                *testing.T
	ctx              context.Context
	waitReceive      *time.Duration
	msgReceiptHandle *string
	errRec           error
	errChg           error
	errDel           error
	msgOut           *sqs.ReceiveMessageOutput
}

func (g mockSQSGateway) ReceiveMessage(ctx context.Context) (*sqs.ReceiveMessageOutput, error) {
	assert.Equal(g.t, g.ctx, ctx, "Unexpected context")
	if g.waitReceive != nil {
		time.Sleep(*g.waitReceive)
	}
	return g.msgOut, g.errRec
}
func (g mockSQSGateway) ChangeMessageVisibility(ctx context.Context, msgReceiptHandle *string) error {
	assert.Equal(g.t, g.ctx, ctx, "Unexpected context")
	return g.errChg
}
func (g mockSQSGateway) DeleteMessage(ctx context.Context, msgReceiptHandle *string) error {
	assert.Equal(g.t, g.ctx, ctx, "Unexpected context")
	return g.errDel
}

var ctx context.Context
var msgId string
var prodQuot model.ProductQuotation
var quotReq api.QuotationRequest
var quotRes *api.QuotationResponse
var quotEnt entity.QuotationEntity
var sqsMsgOut sqs.ReceiveMessageOutput
var snsMsg message.UserPricesMessage
var err error

func initVariables() {

	ctx = context.TODO()
	msgId = uuid.New().String()
	prodQuot = model.ProductQuotation{
		ProductId:     "A",
		Quantity:      1,
		ProductName:   "Mock Product A",
		OriginalValue: 100.0,
		Discount:      10.0,
		FinalValue:    90.0,
	}

	quotReq = api.QuotationRequest{
		RequestId: uuid.New().String(),
		UserId:    "ABC",
		Products: []model.QuotationItem{
			{
				ProductId: prodQuot.ProductId,
				Quantity:  prodQuot.Quantity,
			},
		},
	}

	quotEnt = entity.QuotationEntity{
		RequestId:   quotReq.RequestId,
		UserId:      quotReq.UserId,
		ProductList: quotReq.Products,
	}

	snsMsg = message.UserPricesMessage{
		RequestId: quotReq.RequestId,
		UserId:    quotReq.UserId,
	}
	sqsBody, _ := json.Marshal(message.QuotationMessage{
		RequestId: quotReq.RequestId,
		UserId:    quotReq.UserId,
		Products: []model.ProductQuotation{
			prodQuot,
		},
	})
	sqsMsgOut = sqs.ReceiveMessageOutput{
		Messages: []types.Message{
			{
				ReceiptHandle: &msgId,
				MessageAttributes: map[string]types.MessageAttributeValue{
					"RequestId": {StringValue: &quotReq.RequestId},
				},
				Body: aws.String(string(sqsBody)),
			},
		},
	}

	quotRes = &api.QuotationResponse{
		RequestId: quotReq.RequestId,
		UserId:    quotReq.UserId,
		Products: []model.ProductQuotation{
			prodQuot,
		},
	}
	err = errors.New("Mock Error")
}
