package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/lib_common/model"
	"github.com/ricardo-comar/organic-cache/quotation_provider/gateway"
)

func main() {
	lambdaHandler := awsHandler{dynGtw: gateway.NewDynamoGateway(), sqsGtw: gateway.NewSQSGateway()}
	lambda.Start(lambdaHandler.handleMessages)
}

type awsHandler struct {
	dynGtw gateway.DynamoGateway
	sqsGtw gateway.SQSGateway
}

func (l awsHandler) handleMessages(ctx context.Context, snsEvent events.SNSEvent) error {

	inicioProc := time.Now()

	for _, record := range snsEvent.Records {
		inicioMsg := time.Now()

		err := handleMessage(l, ctx, record)
		if err != nil {
			log.Printf("Erro processando mensagem: %+v", err)
		}

		log.Printf("Finalizando - mensagem %s em %dms", record.SNS.MessageID, time.Since(inicioMsg).Milliseconds())
	}

	log.Printf("Finalizando - processamento em %dms", time.Since(inicioProc).Milliseconds())
	return nil
}

func handleMessage(l awsHandler, ctx context.Context, event events.SNSEventRecord) error {

	log.Printf("Processando mensagem: %+v", event)

	msg := message.UserPricesMessage{}
	if err := json.Unmarshal([]byte(event.SNS.Message), &msg); err != nil {
		log.Printf("Erro transformando mensagem: %+v - %+v", err, event.SNS.Message)
		return err
	}

	productPrices, err := l.dynGtw.QueryProductPrice(msg.UserId)
	if err != nil {
		log.Printf("Erro buscando por produtos calculados: %+v", err)
		return err
	}

	if productPrices == nil {

		log.Println("Nenhuma cotação encontrada, solicitando tabela de preços")
		l.sqsGtw.RecalcMessage(ctx, &msg)

	} else {

		var quotationResponse = &message.QuotationMessage{}
		quotationResponse.RequestId = msg.RequestId
		quotationResponse.UserId = msg.UserId
		quotationResponse.Products = []model.ProductQuotation{}

		quotationRequest, err := l.dynGtw.QueryRequest(msg.RequestId)

		if err != nil {
			log.Printf("Erro ao recuperar o quotation request em base : %+v", err)

		} else {

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
		l.sqsGtw.NotifyQuotation(ctx, *quotationResponse)
	}

	return err

}
