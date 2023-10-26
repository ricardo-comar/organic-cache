package main

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/stretchr/testify/assert"
)

var ctx context.Context
var sendMessageCalled, queryUsersCalled int
var usersDyn []entity.UserEntity
var msgId = uuid.New()
var err error

func initVariables() {

	sendMessageCalled = 0
	queryUsersCalled = 0

	usersDyn = []entity.UserEntity{
		{UserId: "AAA",
			TTL: "now",
		},
		{UserId: "BBB",
			TTL: "now",
		},
	}

	err = errors.New("Mock Error")

	ctx = context.TODO()

}
func TestSuccess(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t, users: usersDyn},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx, messageId: &msgId},
	}

	reqEvent := events.CloudWatchEvent{}
	respEvent, respErr := lambdaHandler.eventHandler(ctx, reqEvent)

	assert.Nil(t, respErr, "Unexpected error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(respEvent, reqEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 1, queryUsersCalled, "Unexpected QueryUsers calls")
	assert.Equal(t, 2, sendMessageCalled, "Unexpected SendMessage calls")

}
func TestErrorQueryUsers(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t, err: err},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx},
	}

	reqEvent := events.CloudWatchEvent{}
	respEvent, respErr := lambdaHandler.eventHandler(ctx, reqEvent)

	assert.Equal(t, err, respErr, "Unexpected error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(respEvent, reqEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 1, queryUsersCalled, "Unexpected QueryUsers calls")
	assert.Equal(t, 0, sendMessageCalled, "Unexpected SendMessage calls")

}
func TestEmptyQueryUsers(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t, users: []entity.UserEntity{}},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx},
	}

	reqEvent := events.CloudWatchEvent{}
	respEvent, respErr := lambdaHandler.eventHandler(ctx, reqEvent)

	assert.Nil(t, respErr, "Unexpected error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(respEvent, reqEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 1, queryUsersCalled, "Unexpected QueryUsers calls")
	assert.Equal(t, 0, sendMessageCalled, "Unexpected SendMessage calls")

}
func TestOneQueryUsers(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t, users: usersDyn[0:1]},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx, messageId: &msgId},
	}

	reqEvent := events.CloudWatchEvent{}
	respEvent, respErr := lambdaHandler.eventHandler(ctx, reqEvent)

	assert.Nil(t, respErr, "Unexpected error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(respEvent, reqEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 1, queryUsersCalled, "Unexpected QueryUsers calls")
	assert.Equal(t, 1, sendMessageCalled, "Unexpected SendMessage calls")

}
func TestErrorErrorSendMessage(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t, users: usersDyn},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx, err: err, messageId: &msgId},
	}

	reqEvent := events.CloudWatchEvent{}
	respEvent, respErr := lambdaHandler.eventHandler(ctx, reqEvent)

	assert.Equal(t, err, respErr, "Unexpected error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(respEvent, reqEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 1, queryUsersCalled, "Unexpected QueryUsers calls")
	assert.Equal(t, 1, sendMessageCalled, "Unexpected SendMessage calls")

}

type mockSQSGateway struct {
	t         *testing.T
	ctx       context.Context
	user      *message.UserMessage
	messageId *uuid.UUID
	err       error
}

func (g mockSQSGateway) SendMessage(ctx context.Context, user *message.UserMessage) (*string, error) {
	sendMessageCalled += 1
	assert.Equal(g.t, g.ctx, ctx, "Unexpected ctx")

	if g.user == nil {
		assert.Nil(g.t, g.user, "User is not nil")

	} else {
		if diff := deep.Equal(g.user, user); diff != nil {
			g.t.Error("Invalid User: ", diff)
		}
	}

	if g.messageId != nil {
		messageId := g.messageId.String()
		return &messageId, g.err
	}

	return nil, g.err
}

type mockDynamoGateway struct {
	t     *testing.T
	users []entity.UserEntity
	err   error
}

func (g mockDynamoGateway) QueryUsers() ([]entity.UserEntity, error) {
	queryUsersCalled += 1
	return g.users, g.err

}
