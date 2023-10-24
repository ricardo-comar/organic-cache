package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/message"

	"github.com/ricardo-comar/organic-cache/price_calc/gateway"
	"github.com/ricardo-comar/organic-cache/price_calc/service"
)

var cfg aws.Config
var dyncli dynamodb.Client
var snscli sns.Client

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
	lambda.Start(handleMessages)
}

func handleMessages(ctx context.Context, sqsEvent events.SQSEvent) error {

	inicioProc := time.Now()

	for _, record := range sqsEvent.Records {
		inicioMsg := time.Now()

		log.Printf("Processando mensagem: %s", record.Body)
		user := entity.UserEntity{}
		json.Unmarshal([]byte(record.Body), &user)

		err := service.GenerateUserPrices(&dyncli, &user)

		if err != nil {

			log.Printf("Error handling message: %+v", err)

		} else {

			if requestId, found := record.MessageAttributes["RequestId"]; found {
				gateway.NotifyQuotation(ctx, &snscli, message.UserPricesMessage{
					UserId: user.UserId, RequestId: *requestId.StringValue,
				})
			}
		}

		log.Printf("Finalizando - mensagem %s em %dms", record.MessageId, time.Since(inicioMsg).Milliseconds())
	}

	log.Printf("Finalizando - processamento em %dms", time.Since(inicioProc).Milliseconds())
	return nil
}
