package main

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/stretchr/testify/assert"
)

var ctx context.Context
var calledGenerateUserPrices int
var calledNotifyQuotation int
var userMsg *message.UserMessage
var snsMsg message.UserPricesMessage
var msgId = uuid.New().String()
var sqsEvent events.SQSEvent

func initVariables() {
	calledGenerateUserPrices = 0
	calledNotifyQuotation = 0

	ctx = context.TODO()

	userMsg = &message.UserMessage{
		UserId: "ABC",
	}
	userMsgBody, _ := json.Marshal(userMsg)

	snsMsg = message.UserPricesMessage{
		UserId:    userMsg.UserId,
		RequestId: uuid.New().String(),
	}

	sqsEvent = events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:     uuid.New().String(),
				ReceiptHandle: uuid.New().String(),
				Body:          string(userMsgBody),
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"RequestId": {StringValue: &snsMsg.RequestId},
				},
			},
		},
	}

}
func TestSuccess(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		ps: mockPricesService{t: t, userMsg: userMsg},
		sg: mockSNSGateway{t: t, ctx: ctx, msg: snsMsg, messageId: &msgId},
	}

	respErr := lambdaHandler.handleMessages(ctx, sqsEvent)

	assert.Nil(t, respErr, "respErr not nil")
	assert.Equal(t, 1, calledGenerateUserPrices, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 1, calledNotifyQuotation, "Unexpected calledNotifyQuotation calls")

}
func TestSuccessNoRequestId(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		ps: mockPricesService{t: t, userMsg: userMsg},
		sg: mockSNSGateway{t: t, ctx: ctx, msg: snsMsg, messageId: &msgId},
	}

	sqsEvent.Records[0].MessageAttributes = make(map[string]events.SQSMessageAttribute)
	respErr := lambdaHandler.handleMessages(ctx, sqsEvent)

	assert.Nil(t, respErr, "respErr not nil")
	assert.Equal(t, 1, calledGenerateUserPrices, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 0, calledNotifyQuotation, "Unexpected calledNotifyQuotation calls")

}
func TestError(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		ps: mockPricesService{t: t, userMsg: userMsg, errGenerateUserPrices: errors.New("mockError")},
		sg: mockSNSGateway{t: t, ctx: ctx, msg: snsMsg, messageId: &msgId},
	}

	respErr := lambdaHandler.handleMessages(ctx, sqsEvent)

	assert.Nil(t, respErr, "respErr not nil")
	assert.Equal(t, 1, calledGenerateUserPrices, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 0, calledNotifyQuotation, "Unexpected calledNotifyQuotation calls")

}

type mockPricesService struct {
	t                     *testing.T
	userMsg               *message.UserMessage
	errGenerateUserPrices error
}

func (rs mockPricesService) GenerateUserPrices(user *message.UserMessage) error {
	calledGenerateUserPrices++

	if rs.userMsg == nil {
		assert.Nil(rs.t, user, "UserMessage is not nil")
	} else if diff := deep.Equal(rs.userMsg, user); diff != nil {
		rs.t.Error("Invalid UserMessage: ", diff)
	}
	return rs.errGenerateUserPrices
}

type mockSNSGateway struct {
	t         *testing.T
	ctx       context.Context
	msg       message.UserPricesMessage
	messageId *string
	err       error
}

func (g mockSNSGateway) NotifyQuotation(ctx context.Context, msg message.UserPricesMessage) (*string, error) {
	calledNotifyQuotation++

	assert.Equal(g.t, g.ctx, ctx, "Unexpected context")
	assert.Equal(g.t, g.msg.RequestId, msg.RequestId, "Invalid RequestId")
	assert.Equal(g.t, g.msg.UserId, msg.UserId, "Invalid UserId")

	return g.messageId, g.err
}
