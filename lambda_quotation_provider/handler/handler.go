package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/quotation_provider/gateway"
	"github.com/ricardo-comar/organic-cache/quotation_provider/service"
)

func main() {
	lambdaHandler := awsHandler{service.NewQuotationService(gateway.NewDynamoGateway(), gateway.NewSQSGateway())}
	lambda.Start(lambdaHandler.handleMessages)
}

type awsHandler struct {
	service service.QuotationService
}

func (l awsHandler) handleMessages(ctx context.Context, snsEvent events.SNSEvent) error {

	inicioProc := time.Now()

	for _, record := range snsEvent.Records {
		inicioMsg := time.Now()

		log.Printf("Processando mensagem: %+v", record)
		msg := &message.UserPricesMessage{}
		decoder := json.NewDecoder(strings.NewReader(record.SNS.Message))
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(msg); err != nil {
			log.Printf("Erro transformando mensagem: %+v - %+v", err, record.SNS.Message)
			continue
		}

		err := l.service.GenerateUserQuotation(ctx, msg)
		if err != nil {
			log.Printf("Erro processando mensagem: %+v", err)
		}

		log.Printf("Finalizando - mensagem %s em %dms", record.SNS.MessageID, time.Since(inicioMsg).Milliseconds())
	}

	log.Printf("Finalizando - processamento em %dms", time.Since(inicioProc).Milliseconds())
	return nil
}
