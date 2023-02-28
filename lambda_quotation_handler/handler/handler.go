package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/model"
)

var cfg aws.Config
var snscli sns.Client
var dyncli dynamodb.Client

func init() {
	cfg, _ = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("AWS_REGION")
		return nil
	})

	dyncli = *dynamodb.NewFromConfig(cfg)
	snscli = *sns.NewFromConfig(cfg)

	localendpoint, found := os.LookupEnv("LOCALSTACK_HOSTNAME")
	if found {
		dyncli = *dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL("http://"+localendpoint+":4566")))
		snscli = *sns.New(sns.Options{Credentials: cfg.Credentials, EndpointResolver: sns.EndpointResolverFromURL("http://" + localendpoint + ":" + os.Getenv("EDGE_PORT"))})
	}
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	var response model.SocketResponse
	switch request.RequestContext.RouteKey {

	case "$connect":
		response = model.SocketResponse{Message: "connected !"}

	case "$disconnect":
		response = model.SocketResponse{Message: "disconnected !"}

	case "PING":
		response = model.SocketResponse{Message: "pong !"}

	case "MESSAGE":

		msg := model.SocketRequest{}
		err := json.Unmarshal([]byte(request.Body), &msg)
		if err != nil {
			log.Printf("Error parsing quotation request: %+v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
		}
		log.Printf("Message: %+v", msg)

		quotationReq := model.QuotationRequest{}
		quotationReq.UserId = msg.Payload.Message.UserId
		quotationReq.ProductList = msg.Payload.Message.ProductList

		quotationReq.RequestId = uuid.New().String()
		quotationReq.TTL = strconv.FormatInt(time.Now().Add(time.Minute*5).UnixNano(), 10)
		quotationReq.ConnectionId = request.RequestContext.ConnectionID

		log.Printf("Sending quotation: %+v", quotationReq)
		err = gateway.SaveQuotationRequest(&dyncli, &quotationReq)
		if err != nil {
			log.Printf("Error saving quotation request: %+v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}

		_, err = gateway.NotifyQuotation(ctx, &snscli, model.MessageEntity{
			RequestId: quotationReq.RequestId,
			UserId:    quotationReq.UserId,
		})
		if err != nil {
			log.Printf("Error notifying quotation: %+v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}

		response = model.SocketResponse{Message: "quotation under analisys"}
	default:
		response = model.SocketResponse{Message: "route not implemented"}
	}

	resp, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{Body: string(resp), StatusCode: http.StatusOK}, nil
}
