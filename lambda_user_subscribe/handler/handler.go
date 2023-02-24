package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/service"
)

var cfg aws.Config
var dyncli dynamodb.Client
var sqscli sqs.Client

func init() {
	cfg, _ = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("AWS_REGION")
		return nil
	})

	dyncli = *dynamodb.NewFromConfig(cfg)
	sqscli = *sqs.NewFromConfig(cfg)

	localendpoint, found := os.LookupEnv("LOCALSTACK_HOSTNAME")
	if found {
		dyncli = *dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL("http://"+localendpoint+":4566")))
		sqscli = *sqs.New(sqs.Options{Credentials: cfg.Credentials, EndpointResolver: sqs.EndpointResolverFromURL("http://" + localendpoint + ":" + os.Getenv("EDGE_PORT"))})
	}
}

func main() {
	lambda.Start(handleRequest)
}

type MyEvent struct {
	Records []MyEventRecord `json:"Records"`

	HTTPMethod string `json:"httpMethod"`
	Body       string `json:"body"`
}

type MyEventRecord struct {
	SNS MyEntity `json:"Sns"`
}

type MyEntity struct {
	Message string `json:"Message"`
}

func handleRequest(ctx context.Context, request MyEvent) (events.APIGatewayProxyResponse, error) {

	var body string
	if len(request.Records) > 0 {
		body = request.Records[0].SNS.Message
	} else {
		body = request.Body
	}

	entity, err := service.CreateEntity(body)
	if err != nil {
		log.Println("Invalid content: ", request.Body)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	err = gateway.SaveActiveUser(&dyncli, entity)
	if err != nil {
		log.Println("Error saving subscription: ", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusCreated}, err

}
