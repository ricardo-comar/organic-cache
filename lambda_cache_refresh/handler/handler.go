package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ricardo-comar/organic-cache/gateway"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var cfg aws.Config

func init() {
	cfg, _ = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("AWS_REGION")
		return nil
	})
}

func main() {
	lambda.Start(eventHandler)
}

func eventHandler(ctx context.Context, event events.CloudWatchEvent) (events.CloudWatchEvent, error) {

	log.Printf("Iniciando busca por usuários subscritos")
	inicio := time.Now()

	users, error := gateway.QueryUsers(cfg)
	log.Printf("%d usuários encontrados", len(users))

	for _, user := range users {
		msgId, _ := gateway.SendMessage(ctx, cfg, user)
		log.Printf("Usuário %s enviado: %s", user.ID, *msgId)
	}

	log.Printf("Finalizando com %d usuários em %dms", len(users), time.Now().Sub(inicio).Milliseconds())

	return event, error

}
