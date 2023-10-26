package service

import (
	"context"
	"log"

	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/lib_common/model"
	"github.com/ricardo-comar/organic-cache/quotation_provider/gateway"
)

type quotationService struct {
	dynGtw gateway.DynamoGateway
	sqsGtw gateway.SQSGateway
}

func NewQuotationService(dg gateway.DynamoGateway, sqsGtw gateway.SQSGateway) QuotationService {
	rs := &quotationService{
		dynGtw: dg,
		sqsGtw: sqsGtw,
	}
	return QuotationService(rs)
}

type QuotationService interface {
	GenerateUserQuotation(ctx context.Context, user *message.UserPricesMessage) error
}

func (s quotationService) GenerateUserQuotation(ctx context.Context, msg *message.UserPricesMessage) error {

	productPrices, err := s.dynGtw.QueryProductPrice(msg.UserId)
	if err != nil {
		log.Printf("Erro buscando por produtos calculados: %+v", err)
		return err
	}

	if productPrices == nil {

		log.Println("Nenhuma cotação encontrada, solicitando tabela de preços")
		msgId, err := s.sqsGtw.RecalcMessage(ctx, msg)
		log.Println("Recalculo enviado: ", msgId)

		if err != nil {
			log.Printf("Erro enviando mensagem de recálculo: %+v", err)
			return err
		}

	} else {

		var quotationResponse = &message.QuotationMessage{}
		quotationResponse.RequestId = msg.RequestId
		quotationResponse.UserId = msg.UserId
		quotationResponse.Products = []model.ProductQuotation{}

		quotationRequest, err := s.dynGtw.QueryRequest(msg.RequestId)

		if err != nil {
			log.Printf("Erro ao recuperar o quotation request em base : %+v", err)
			return err
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
		quotId, err := s.sqsGtw.NotifyQuotation(ctx, quotationResponse)
		log.Printf("Cotação notificada: %+v", quotId)

		if err != nil {
			log.Printf("Erro enviando notificação de cotação: %+v", err)
			return err
		}

	}

	return nil

}
