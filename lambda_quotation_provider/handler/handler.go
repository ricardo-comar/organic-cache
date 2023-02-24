package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	inv "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/model"
)

var cfg aws.Config
var dyncli dynamodb.Client
var sqscli sqs.Client
var lamcli lambda.Client

func init() {
	cfg, _ = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("AWS_REGION")
		return nil
	})

	dyncli = *dynamodb.NewFromConfig(cfg)
	sqscli = *sqs.NewFromConfig(cfg)
	lamcli = *lambda.NewFromConfig(cfg)

	localendpoint, found := os.LookupEnv("LOCALSTACK_HOSTNAME")
	if found {
		localhost := "http://" + localendpoint + ":" + os.Getenv("EDGE_PORT")
		dyncli = *dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL(localhost)))
		sqscli = *sqs.New(sqs.Options{Credentials: cfg.Credentials, EndpointResolver: sqs.EndpointResolverFromURL(localhost)})
		lamcli = *lambda.New(lambda.Options{Credentials: cfg.Credentials, EndpointResolver: lambda.EndpointResolverFromURL(localhost)})
	}
}

func main() {
	inv.Start(handleMessages)
}

func handleMessages(ctx context.Context, sqsEvent events.SQSEvent) error {

	inicioProc := time.Now()

	for _, record := range sqsEvent.Records {
		inicioMsg := time.Now()

		err := handleMessage(ctx, record)
		if err != nil {
			log.Printf("Erro processando mensagem: %+v", err)
		}

		log.Printf("Finalizando - mensagem %s em %dms", record.MessageId, time.Now().Sub(inicioMsg).Milliseconds())
	}

	log.Printf("Finalizando - processamento em %dms", time.Now().Sub(inicioProc).Milliseconds())
	return nil
}

func handleMessage(ctx context.Context, msg events.SQSMessage) error {

	log.Printf("Processando mensagem: %+v", msg)
	msgId := msg.MessageAttributes["RequestId"].StringValue

	var quotationRequest model.QuotationRequest
	json.Unmarshal([]byte(msg.Body), &quotationRequest)

	productPrices, err := gateway.QueryProductPrice(dyncli, quotationRequest.UserId)
	if err != nil {
		log.Printf("Erro buscando por produtos calculados: %+v", err)
		return err
	}

	if productPrices == nil {

		log.Println("Nenhuma cotação encontrada, postergando processamento")

		retries := "1"
		if retryCount := msg.MessageAttributes["RetryCount"].StringValue; retryCount != nil {
			rCount, _ := strconv.Atoi(*retryCount)
			retries = strconv.Itoa(rCount + 1)
			log.Printf("Retries: %v | %v | %v", *retryCount, rCount, retries)
			if rCount >= 10 {
				return nil
			}

		} else {
			log.Println("Solicitando cálculo de preços")
			gateway.RecalcMessage(ctx, &sqscli, quotationRequest.UserId)
		}

		time.Sleep(time.Second)
		gateway.RetryMessage(ctx, &sqscli, msgId, &retries, quotationRequest)

	} else {

		var quotationResponse *model.QuotationEntity
		quotationResponse = &model.QuotationEntity{}
		quotationResponse.Id = *msgId
		quotationResponse.Products = []model.ProductQuotation{}

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

		log.Printf("Cotação realizada: %+v", *quotationResponse)
		gateway.SendResponse(&ctx, quotationResponse)
	}

	return err

}
