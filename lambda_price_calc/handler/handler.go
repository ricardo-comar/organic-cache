package main

import (
	"context"
	"log"
	"os"
	"time"

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
	lambda.Start(handleMessages)
}

func handleMessages(ctx context.Context, sqsEvent events.SQSEvent) error {

	inicioProc := time.Now()

	for _, message := range sqsEvent.Records {
		inicioMsg := time.Now()

		handleMessage(message.Body)

		log.Printf("Finalizando - mensagem %s em %dms", message.MessageId, time.Now().Sub(inicioMsg).Milliseconds())
	}

	log.Printf("Finalizando - processamento em %dms", time.Now().Sub(inicioProc).Milliseconds())
	return nil
}

func handleMessage(msg string) (interface{}, error) {

	log.Printf("Processando mensagem: %s", msg)

	return nil, nil

}
