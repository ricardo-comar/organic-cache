package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ricardo-comar/organic-cache/user_subscribe/gateway"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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
	lambda.Start(eventHandler)
}

func eventHandler(ctx context.Context, event events.CloudWatchEvent) (events.CloudWatchEvent, error) {

	log.Printf("Iniciando busca por usuários subscritos")
	inicio := time.Now()

	users, error := gateway.QueryUsers(&dyncli)
	log.Printf("%d usuários encontrados", len(users))

	for _, user := range users {
		msgId, _ := gateway.SendMessage(ctx, &sqscli, user)
		log.Printf("Usuário %s enviado: %s", user.ID, *msgId)
	}

	log.Printf("Finalizando com %d usuários em %dms", len(users), time.Now().Sub(inicio).Milliseconds())

	return event, error

}
