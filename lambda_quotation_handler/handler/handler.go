package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/model"
	"github.com/ricardo-comar/organic-cache/service"
)

var cfg aws.Config
var dyncli dynamodb.Client
var sqscli sqs.Client
var snscli sns.Client

func init() {
	cfg, _ = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("AWS_REGION")
		return nil
	})

	dyncli = *dynamodb.NewFromConfig(cfg)
	sqscli = *sqs.NewFromConfig(cfg)
	snscli = *sns.NewFromConfig(cfg)

	localendpoint, found := os.LookupEnv("LOCALSTACK_HOSTNAME")
	if found {
		dyncli = *dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL("http://"+localendpoint+":4566")))
		sqscli = *sqs.New(sqs.Options{Credentials: cfg.Credentials, EndpointResolver: sqs.EndpointResolverFromURL("http://" + localendpoint + ":" + os.Getenv("EDGE_PORT"))})
		snscli = *sns.New(sns.Options{Credentials: cfg.Credentials, EndpointResolver: sns.EndpointResolverFromURL("http://" + localendpoint + ":" + os.Getenv("EDGE_PORT"))})
	}
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod != http.MethodPost {
		return events.APIGatewayProxyResponse{Body: http.StatusText(http.StatusMethodNotAllowed), StatusCode: http.StatusMethodNotAllowed}, nil
	}

	quotationReq := model.QuotationRequest{}
	err := json.Unmarshal([]byte(request.Body), &quotationReq)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	topicKey, err := gateway.SubscribeUser(ctx, &snscli, gateway.UserMessage{ID: quotationReq.UserId})
	log.Printf("Usuário %v notificado: %v", quotationReq.UserId, topicKey)

	reqId := uuid.New().String()

	_, err = gateway.SendMessage(ctx, &sqscli, quotationReq, &reqId)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	quotationResponse, err := service.ResponseWait(&ctx, &dyncli, reqId)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	if quotationResponse == nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusRequestTimeout}, err
	}

	response, _ := json.Marshal(quotationResponse)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(response)}, err

}
