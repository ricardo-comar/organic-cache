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
	"github.com/ricardo-comar/organic-cache/lib_common/api"
	"github.com/ricardo-comar/organic-cache/quotation_handler/service"
)

var cfg aws.Config
var snscli *sns.Client
var sqscli *sqs.Client
var dyncli *dynamodb.Client

func init() {
	cfg, _ = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("AWS_REGION")
		return nil
	})

	dyncli = dynamodb.NewFromConfig(cfg)
	snscli = sns.NewFromConfig(cfg)
	sqscli = sqs.NewFromConfig(cfg)

	localendpoint, found := os.LookupEnv("LOCALSTACK_HOSTNAME")
	if found {
		dyncli = dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL("http://"+localendpoint+":4566")))
		snscli = sns.New(sns.Options{Credentials: cfg.Credentials, EndpointResolver: sns.EndpointResolverFromURL("http://" + localendpoint + ":" + os.Getenv("EDGE_PORT"))})
		sqscli = sqs.New(sqs.Options{Credentials: cfg.Credentials, EndpointResolver: sqs.EndpointResolverFromURL("http://" + localendpoint + ":" + os.Getenv("EDGE_PORT"))})
	}
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	req := api.QuotationRequest{}
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error parsing request: %+v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}
	log.Printf("Message: %+v", req)

	err = service.RequestQuotation(ctx, snscli, dyncli, req, request.RequestContext.RequestID)
	if err != nil {
		log.Printf("Error saving quotation request: %+v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	response, err := service.WaitForResponse(ctx, sqscli, request.RequestContext.RequestID)

	if err == nil {
		log.Printf("Response: %+v", response)
		resp, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{Body: string(resp), StatusCode: http.StatusOK}, nil
	}

	switch err.(type) {
	case *service.TimeoutError:
		log.Printf("Timeout waiting for quotation response: %+v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusRequestTimeout}, err
	}

	log.Printf("Error waiting for quotation response: %+v", err)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
}
