package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/model"
	"github.com/ricardo-comar/organic-cache/service"
)

var cfg aws.Config
var sqscli sqs.Client
var snscli sns.Client

func init() {
	cfg, _ = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("AWS_REGION")
		return nil
	})

	sqscli = *sqs.NewFromConfig(cfg)
	snscli = *sns.NewFromConfig(cfg)

	localendpoint, found := os.LookupEnv("LOCALSTACK_HOSTNAME")
	if found {
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

	reqId := uuid.New().String()

	_, err = gateway.SendMessage(ctx, &sqscli, quotationReq, &reqId)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)

	defer cancel()

	quotationResponse, err := service.ResponseWait(ctx, &snscli, &reqId)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: quotationResponse}, err

}
