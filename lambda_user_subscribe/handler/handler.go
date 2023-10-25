package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/user_subscribe/gateway"
)

func main() {
	lambdaHandler := awsHandler{dynGtw: gateway.NewDynamoGateway(), sqsGtw: gateway.NewSQSGateway()}
	lambda.Start(lambdaHandler.handleRequest)
}

type awsHandler struct {
	dynGtw gateway.DynamoGateway
	sqsGtw gateway.SQSGateway
}

func (l awsHandler) handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	entity, err := entity.NewUserEntity(request.Body)
	if err != nil {
		log.Println("Invalid content: ", request.Body)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	userSub, err := l.dynGtw.QuerySubscription(entity.UserId)
	if err == nil && userSub == nil {
		log.Println("New subscription, asking for price recalculation: ", entity.UserId)
		l.sqsGtw.SendMessage(ctx, &message.UserMessage{UserId: entity.UserId})
	}

	err = l.dynGtw.SaveActiveUser(entity)
	if err != nil {
		log.Println("Error saving subscription: ", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusCreated}, err

}
