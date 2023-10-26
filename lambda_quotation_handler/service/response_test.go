package service_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/go-test/deep"
	"github.com/ricardo-comar/organic-cache/quotation_handler/service"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	initVariables()

	respService := service.NewResponseService(time.Second, mockSQSGateway{t: t, ctx: ctx, msgReceiptHandle: &msgId, msgOut: &sqsMsgOut})

	respQuot, respErr := respService.WaitForResponse(ctx, quotReq.RequestId)
	assert.Nil(t, respErr, "Unexpected Error")
	if diff := deep.Equal(quotRes, respQuot); diff != nil {
		t.Error("Unexpected response : ", diff)
	}
}

func TestErrorReceiveMessage(t *testing.T) {
	initVariables()

	respService := service.NewResponseService(time.Second, mockSQSGateway{t: t, ctx: ctx, errRec: err})

	respQuot, respErr := respService.WaitForResponse(ctx, quotReq.RequestId)
	assert.Nil(t, respQuot, "Unexpected response")
	assert.Equal(t, err, respErr, "Unexpected error")
}

func TestErrorMessageBody(t *testing.T) {
	initVariables()

	sqsMsgOut.Messages[0].Body = aws.String("ABC")

	respService := service.NewResponseService(time.Second, mockSQSGateway{t: t, ctx: ctx, msgReceiptHandle: &msgId, msgOut: &sqsMsgOut})

	respQuot, respErr := respService.WaitForResponse(ctx, quotReq.RequestId)
	assert.Nil(t, respQuot, "Unexpected response")
	assert.NotNil(t, respErr, "Empty error")
}

func TestErrorDeleteMessage(t *testing.T) {
	initVariables()

	respService := service.NewResponseService(time.Second, mockSQSGateway{t: t, ctx: ctx, errDel: err, msgReceiptHandle: &msgId, msgOut: &sqsMsgOut})

	respQuot, respErr := respService.WaitForResponse(ctx, quotReq.RequestId)
	assert.Nil(t, respQuot, "Unexpected response")
	assert.NotNil(t, respErr, "Empty error")
}

func TestErrorMessageVisibility(t *testing.T) {
	initVariables()

	sqsMsgOut.Messages[0].MessageAttributes = make(map[string]types.MessageAttributeValue)
	respService := service.NewResponseService(time.Second, mockSQSGateway{t: t, ctx: ctx, errChg: err, msgReceiptHandle: &msgId, msgOut: &sqsMsgOut})

	respQuot, respErr := respService.WaitForResponse(ctx, quotReq.RequestId)
	assert.Nil(t, respQuot, "Unexpected response")
	assert.NotNil(t, respErr, "Empty error")
}

func TestErrorTimeout(t *testing.T) {
	initVariables()

	d := time.Second * 2
	respService := service.NewResponseService(time.Second/10, mockSQSGateway{t: t, ctx: ctx, waitReceive: &d})

	respQuot, respErr := respService.WaitForResponse(ctx, quotReq.RequestId)
	assert.Nil(t, respQuot, "Unexpected response")
	assert.Equal(t, &service.TimeoutError{}, respErr, "Empty error")
}
