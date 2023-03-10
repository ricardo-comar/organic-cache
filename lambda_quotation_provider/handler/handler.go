package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	inv "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/model"
)

var cfg aws.Config
var dyncli dynamodb.Client
var sqscli sqs.Client
var gtwcli apigatewaymanagementapi.Client

func init() {
	cfg, _ = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("AWS_REGION")
		return nil
	})

	// apiPath := os.Getenv("API_PATH")

	dyncli = *dynamodb.NewFromConfig(cfg)
	sqscli = *sqs.NewFromConfig(cfg)
	gtwcli = *apigatewaymanagementapi.NewFromConfig(cfg)

	localendpoint, found := os.LookupEnv("LOCALSTACK_HOSTNAME")
	if found {
		localhost := "http://" + localendpoint + ":" + os.Getenv("EDGE_PORT")
		dyncli = *dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL(localhost)))
		sqscli = *sqs.New(sqs.Options{Credentials: cfg.Credentials, EndpointResolver: sqs.EndpointResolverFromURL(localhost)})
		gtwcli = *apigatewaymanagementapi.New(apigatewaymanagementapi.Options{Credentials: cfg.Credentials, EndpointResolver: apigatewaymanagementapi.EndpointResolverFromURL(localhost)})
	}
}

func main() {
	inv.Start(handleMessages)
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

	msg := model.MessageEntity{}
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

		log.Println("Nenhuma cota????o encontrada, solicitando tabela de pre??os")
		gateway.RecalcMessage(ctx, &sqscli, &msg)

	} else {

		var quotationResponse *model.QuotationEntity
		quotationResponse = &model.QuotationEntity{}
		quotationResponse.Id = msg.RequestId
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

						log.Printf("Cota????o de produto: %+v", productQuotation)
						quotationResponse.Products = append(quotationResponse.Products, productQuotation)
					}

				}
			}
		}

		log.Printf("Cota????o realizada: %+v", *quotationResponse)
		gateway.SendResponse(&gtwcli, ctx, *quotationResponse, quotationRequest.ConnectionId)
	}

	return err

}
