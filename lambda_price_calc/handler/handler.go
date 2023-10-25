package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardo-comar/organic-cache/lib_common/message"

	"github.com/ricardo-comar/organic-cache/price_calc/gateway"
	"github.com/ricardo-comar/organic-cache/price_calc/service"
)

func main() {
	lambdaHandler := awsHandler{ps: service.NewPricesService(gateway.NewDynamoGateway()), sg: gateway.NewSQSGateway()}
	lambda.Start(lambdaHandler.handleMessages)
}

type awsHandler struct {
	ps service.PricesService
	sg gateway.SNSGateway
}

func (l awsHandler) handleMessages(ctx context.Context, sqsEvent events.SQSEvent) error {

	inicioProc := time.Now()

	for _, record := range sqsEvent.Records {
		inicioMsg := time.Now()

		log.Printf("Processando mensagem: %s", record.Body)
		user := message.UserMessage{}
		json.Unmarshal([]byte(record.Body), &user)

		err := l.ps.GenerateUserPrices(&user)

		if err != nil {

			log.Printf("Error handling message: %+v", err)

		} else {

			if requestId, found := record.MessageAttributes["RequestId"]; found {
				l.sg.NotifyQuotation(ctx, message.UserPricesMessage{
					UserId: user.UserId, RequestId: *requestId.StringValue,
				})
			}
		}

		log.Printf("Finalizando - mensagem %s em %dms", record.MessageId, time.Since(inicioMsg).Milliseconds())
	}

	log.Printf("Finalizando - processamento em %dms", time.Since(inicioProc).Milliseconds())
	return nil
}
