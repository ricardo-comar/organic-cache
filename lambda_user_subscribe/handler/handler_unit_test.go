package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/stretchr/testify/assert"
)

var ctx context.Context
var sendMessageCalled, querySubscriptionCalled, saveActiveUserCalled int
var userDyn entity.UserEntity
var msgId = uuid.New()
var mockError = errors.New("mock_error")
var gatewayRequestContext events.APIGatewayProxyRequestContext
var gatewayRequest events.APIGatewayProxyRequest

func initVariables() {

	sendMessageCalled = 0
	querySubscriptionCalled = 0
	saveActiveUserCalled = 0

	userDyn = entity.UserEntity{UserId: "AAA", TTL: "now"}

	ctx = context.TODO()

	gatewayRequestContext = events.APIGatewayProxyRequestContext{
		RequestID: uuid.New().String(),
	}
	body, _ := json.Marshal(userDyn)

	gatewayRequest = events.APIGatewayProxyRequest{
		RequestContext: gatewayRequestContext,
		Body:           string(body),
	}

}

func TestSuccessNewSub(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t, userId: userDyn.UserId},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx, messageId: &msgId},
	}

	gatExpected := events.APIGatewayProxyResponse{StatusCode: http.StatusAccepted}
	respEvent, respErr := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.Nil(t, respErr, "Unexpected error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(gatExpected, respEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 1, querySubscriptionCalled, "Unexpected QuerySubscription calls")
	assert.Equal(t, 1, sendMessageCalled, "Unexpected SendMessage calls")
	assert.Equal(t, 1, saveActiveUserCalled, "Unexpected SaveActiveUser calls")

}
func TestSuccessExistingUser(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t, userId: userDyn.UserId, user: &userDyn},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx, messageId: &msgId},
	}

	gatExpected := events.APIGatewayProxyResponse{StatusCode: http.StatusAccepted}
	respEvent, respErr := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.Nil(t, respErr, "Unexpected error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(gatExpected, respEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 1, querySubscriptionCalled, "Unexpected QuerySubscription calls")
	assert.Equal(t, 0, sendMessageCalled, "Unexpected SendMessage calls")
	assert.Equal(t, 1, saveActiveUserCalled, "Unexpected SaveActiveUser calls")

}

func TestInvalidUserId(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx},
	}

	gatewayRequest.Body = ""
	gatExpected := events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}
	respEvent, respErr := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.NotNil(t, respErr, " Nil error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(gatExpected, respEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 0, querySubscriptionCalled, "Unexpected QuerySubscription calls")
	assert.Equal(t, 0, sendMessageCalled, "Unexpected SendMessage calls")
	assert.Equal(t, 0, saveActiveUserCalled, "Unexpected SaveActiveUser calls")

}
func TestErrorQuerySubscription(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t, userId: userDyn.UserId, errQuerySubscription: mockError},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx},
	}

	gatExpected := events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	respEvent, respErr := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.Equal(t, mockError, respErr, " Invalid error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(gatExpected, respEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 1, querySubscriptionCalled, "Unexpected QuerySubscription calls")
	assert.Equal(t, 0, sendMessageCalled, "Unexpected SendMessage calls")
	assert.Equal(t, 0, saveActiveUserCalled, "Unexpected SaveActiveUser calls")

}
func TestErrorSendMessage(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t, userId: userDyn.UserId},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx, messageId: &msgId, errSendMessage: mockError},
	}

	gatExpected := events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	respEvent, respErr := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.Equal(t, mockError, respErr, " Invalid error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(gatExpected, respEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 1, querySubscriptionCalled, "Unexpected QuerySubscription calls")
	assert.Equal(t, 1, sendMessageCalled, "Unexpected SendMessage calls")
	assert.Equal(t, 0, saveActiveUserCalled, "Unexpected SaveActiveUser calls")

}
func TestErrorSaveActiveUser(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		dynGtw: mockDynamoGateway{t: t, userId: userDyn.UserId, errSaveActiveUser: mockError},
		sqsGtw: mockSQSGateway{t: t, ctx: ctx, messageId: &msgId},
	}

	gatExpected := events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	respEvent, respErr := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.Equal(t, mockError, respErr, " Invalid error")
	assert.NotNil(t, respEvent, "Nil response")
	if diff := deep.Equal(gatExpected, respEvent); diff != nil {
		t.Error(diff)
	}

	assert.Equal(t, 1, querySubscriptionCalled, "Unexpected QuerySubscription calls")
	assert.Equal(t, 1, sendMessageCalled, "Unexpected SendMessage calls")
	assert.Equal(t, 1, saveActiveUserCalled, "Unexpected SaveActiveUser calls")

}

type mockSQSGateway struct {
	t              *testing.T
	ctx            context.Context
	user           *message.UserMessage
	messageId      *uuid.UUID
	errSendMessage error
}

func (g mockSQSGateway) SendMessage(ctx context.Context, user *message.UserMessage) (*string, error) {
	sendMessageCalled++

	assert.Equal(g.t, g.ctx, ctx, "Unexpected ctx")

	if g.user == nil {
		assert.Nil(g.t, g.user, "User is not nil")

	} else {
		assert.Equal(g.t, g.user.UserId, user.UserId, "Unexpected UserId")
	}

	if g.messageId != nil {
		messageId := g.messageId.String()
		return &messageId, g.errSendMessage
	}

	return nil, g.errSendMessage
}

type mockDynamoGateway struct {
	t                    *testing.T
	userId               string
	user                 *entity.UserEntity
	errQuerySubscription error
	errSaveActiveUser    error
}

func (g mockDynamoGateway) QuerySubscription(userId string) (*entity.UserEntity, error) {
	querySubscriptionCalled++

	assert.Equal(g.t, g.userId, userId, "Unexpected userId")

	return g.user, g.errQuerySubscription

}

func (g mockDynamoGateway) SaveActiveUser(user *entity.UserEntity) error {
	saveActiveUserCalled++
	if g.user == nil {
		assert.Nil(g.t, g.user, "User is not nil")

	} else {
		assert.Equal(g.t, g.user.UserId, user.UserId, "Unexpected UserId")
	}

	return g.errSaveActiveUser
}
