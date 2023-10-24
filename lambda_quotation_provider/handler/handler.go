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
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/lib_common/model"
	"github.com/ricardo-comar/organic-cache/quotation_provider/gateway"
)

var cfg aws.Config
var dyncli dynamodb.Client
var sqscli sqs.Client
var snscli sns.Client

func init() {
	cfg, _ = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("AWS_REGION")
		return nil
	})

	// apiPath := os.Getenv("API_PATH")

	dyncli = *dynamodb.NewFromConfig(cfg)
	sqscli = *sqs.NewFromConfig(cfg)
	snscli = *sns.NewFromConfig(cfg)

	localendpoint, found := os.LookupEnv("LOCALSTACK_HOSTNAME")
	if found {
		localhost := "http://" + localendpoint + ":" + os.Getenv("EDGE_PORT")
		dyncli = *dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL(localhost)))
		sqscli = *sqs.New(sqs.Options{Credentials: cfg.Credentials, EndpointResolver: sqs.EndpointResolverFromURL(localhost)})
		snscli = *sns.New(sns.Options{Credentials: cfg.Credentials, EndpointResolver: sns.EndpointResolverFromURL(localhost)})
	}
}

func main() {
	lambda.Start(handleMessages)
}

func handleMessages(ctx context.Context, snsEvent events.SNSEvent) error {

	inicioProc := time.Now()

	for _, record := range snsEvent.Records {
		inicioMsg := time.Now()

		err := handleMessage(ctx, record)
		if err != nil {
			log.Printf("Erro processando mensagem: %+v", err)
		}

		log.Printf("Finalizando - mensagem %s em %dms", record.SNS.MessageID, time.Now().Sub(inicioMsg).Milliseconds())
	}

	log.Printf("Finalizando - processamento em %dms", time.Now().Sub(inicioProc).Milliseconds())
	return nil
}

func handleMessage(ctx context.Context, event events.SNSEventRecord) error {

	log.Printf("Processando mensagem: %+v", event)

	msg := message.UserPricesMessage{}
	if err := json.Unmarshal([]byte(event.SNS.Message), &msg); err != nil {
		log.Printf("Erro transformando mensagem: %+v - %+v", err, event.SNS.Message)
		return err
	}

	productPrices, err := gateway.QueryProductPrice(dyncli, msg.UserId)
	if err != nil {
		log.Printf("Erro buscando por produtos calculados: %+v", err)
		return err
	}

	if productPrices == nil {

		log.Println("Nenhuma cotação encontrada, solicitando tabela de preços")
		gateway.RecalcMessage(ctx, &sqscli, &msg)

	} else {

		var quotationResponse = &message.QuotationMessage{}
		quotationResponse.RequestId = msg.RequestId
		quotationResponse.UserId = msg.UserId
		quotationResponse.Products = []model.ProductQuotation{}

		quotationRequest, err := gateway.QueryRequest(dyncli, msg.RequestId)

		if err == nil {

			for _, product := range productPrices.Products {
				log.Printf("Produto calculado : %+v", product)

				for _, req := range quotationRequest.ProductList {
					log.Printf("Produto solicitado : %+v", product)

					if req.ProductId == product.ProductId {
						productQuotation := model.ProductQuotation{
							ProductId:     product.ProductId,
							ProductName:   product.ProductName,
							Quantity:      req.Quantity,
							OriginalValue: product.OriginalValue,
							Discount:      product.Discount,
							FinalValue:    (product.Value * req.Quantity),
						}

						log.Printf("Cotação de produto: %+v", productQuotation)
						quotationResponse.Products = append(quotationResponse.Products, productQuotation)
					}

				}
			}
		}

		log.Printf("Cotação realizada: %+v", *quotationResponse)
		gateway.NotifyQuotation(ctx, &sqscli, *quotationResponse)
	}

	return err

}
