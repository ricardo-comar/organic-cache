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
var generateUserQuotationCalled int
var msg *message.UserPricesMessage
var snsEvent events.SNSEvent

func initVariables() {

	generateUserQuotationCalled = 0

	msg = &message.UserPricesMessage{
		RequestId: uuid.New().String(),
		UserId:    "ABC",
	}
	msgBody, _ := json.Marshal(msg)

	snsEvent = events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				SNS: events.SNSEntity{
					MessageID: uuid.New().String(),
					Message:   string(msgBody),
				},
			},
		},
	}
	ctx = context.TODO()

}

func TestSuccess(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		service: mockQuotationService{t: t, ctx: ctx, msg: msg},
	}

	respErr := lambdaHandler.handleMessages(ctx, snsEvent)

	assert.Nil(t, respErr, "Unexpected error")
	assert.Equal(t, 1, generateUserQuotationCalled, "Unexpected GenerateUserQuotation calls")

}
func TestErrorGenerateUserQuotation(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		service: mockQuotationService{t: t, ctx: ctx, msg: msg, errGenerateUserQuotation: errors.New("mock_error")},
	}

	respErr := lambdaHandler.handleMessages(ctx, snsEvent)

	assert.Nil(t, respErr, "Unexpected error")
	assert.Equal(t, 1, generateUserQuotationCalled, "Unexpected GenerateUserQuotation calls")

}
func TestErrorMessageBody(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		service: mockQuotationService{t: t, ctx: ctx},
	}

	snsEvent.Records[0].SNS.Message = "{\"id\":123}"
	respErr := lambdaHandler.handleMessages(ctx, snsEvent)

	assert.Nil(t, respErr, "Unexpected error")
	assert.Equal(t, 0, generateUserQuotationCalled, "Unexpected GenerateUserQuotation calls")

}
func TestSucessAndBodyError(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		service: mockQuotationService{t: t, ctx: ctx, msg: msg},
	}

	snsEvent.Records = append(snsEvent.Records,
		events.SNSEventRecord{
			SNS: events.SNSEntity{
				MessageID: uuid.New().String(),
				Message:   "{\"id\":123}",
			},
		})

	respErr := lambdaHandler.handleMessages(ctx, snsEvent)

	assert.Nil(t, respErr, "Unexpected error")
	assert.Equal(t, 1, generateUserQuotationCalled, "Unexpected GenerateUserQuotation calls")

}

type mockQuotationService struct {
	t                        *testing.T
	ctx                      context.Context
	msg                      *message.UserPricesMessage
	errGenerateUserQuotation error
}

func (s mockQuotationService) GenerateUserQuotation(ctx context.Context, msg *message.UserPricesMessage) error {
	generateUserQuotationCalled++
	assert.Equal(s.t, s.ctx, ctx, "Unexpected ctx")

	if s.msg == nil {
		assert.Nil(s.t, msg, "msg is not nil")

	} else if diff := deep.Equal(s.msg, msg); diff != nil {
		s.t.Error(diff)
	}

	return s.errGenerateUserQuotation
}
