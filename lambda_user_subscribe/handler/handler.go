package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/service"
)

var cfg aws.Config

func init() {
	cfg, _ = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("AWS_REGION")
		return nil
	})
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod != http.MethodPut {
		return events.APIGatewayProxyResponse{Body: http.StatusText(http.StatusMethodNotAllowed), StatusCode: http.StatusMethodNotAllowed}, nil
	}

	entity, err := service.CreateEntity(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	err = gateway.SaveActiveUser(cfg, entity)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusCreated}, err

}
