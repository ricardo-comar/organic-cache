package main

import (
	"context"
	"log"
	"time"

	"github.com/ricardo-comar/organic-cache/cache_refresh/gateway"
	"github.com/ricardo-comar/organic-cache/lib_common/message"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambdaHandler := awsHandler{dynGtw: gateway.NewDynamoGateway(), sqsGtw: gateway.NewSQSGateway()}
	lambda.Start(lambdaHandler.eventHandler)
}

type awsHandler struct {
	dynGtw gateway.DynamoGateway
	sqsGtw gateway.SQSGateway
}

func (l awsHandler) eventHandler(ctx context.Context, event events.CloudWatchEvent) (events.CloudWatchEvent, error) {

	log.Printf("Iniciando busca por usuários subscritos")
	inicio := time.Now()

	users, err := l.dynGtw.QueryUsers()
	if err != nil {
		log.Printf("Erro ao recuperar usuários em base: %+v", err)
		return event, err
	}
	log.Printf("%d usuários encontrados", len(users))

	for _, user := range users {
		msgId, err := l.sqsGtw.SendMessage(ctx, &message.UserMessage{UserId: user.UserId})
		if err != nil {
			log.Printf("Erro ao enviar mensagem via SQS: %+v", err)
			return event, err
		}

		log.Printf("Usuário %s enviado: %s", user.UserId, *msgId)
	}

	log.Printf("Finalizando com %d usuários em %dms", len(users), time.Since(inicio).Milliseconds())

	return event, err

}
