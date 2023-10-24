package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/ricardo-comar/organic-cache/lib_common/api"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/quotation_handler/gateway"
)

func RequestQuotation(ctx context.Context, snscli *sns.Client, dyncli *dynamodb.Client, req api.QuotationRequest, reqId string) error {

	log.Printf("Saving quotation: %+v", req)

	err := gateway.SaveQuotationRequest(dyncli, entity.QuotationEntity{
		RequestId:   reqId,
		UserId:      req.UserId,
		ProductList: req.Products,
		TTL:         strconv.FormatInt(time.Now().Add(time.Minute*5).UnixNano(), 10),
	})
	if err != nil {
		log.Printf("Error saving quotation request: %+v", err)
		return err
	}

	log.Printf("Notifying quotation: %+v", req)
	_, err = gateway.NotifyQuotation(ctx, snscli, message.UserPricesMessage{
		RequestId: reqId,
		UserId:    req.UserId,
	})
	if err != nil {
		log.Printf("Error notifying quotation: %+v", err)
		return err
	}

	return nil
}
