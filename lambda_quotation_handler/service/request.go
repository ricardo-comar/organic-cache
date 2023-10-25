package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/ricardo-comar/organic-cache/lib_common/api"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/quotation_handler/gateway"
)

type requestService struct {
	dg   gateway.DynamoGateway
	snsg gateway.SNSGateway
}

func NewRequestService(dg gateway.DynamoGateway, snsg gateway.SNSGateway) RequestService {
	rs := &requestService{
		dg:   dg,
		snsg: snsg,
	}
	return RequestService(rs)
}

type RequestService interface {
	RequestQuotation(ctx context.Context, req api.QuotationRequest) error
}

func (rs requestService) RequestQuotation(ctx context.Context, req api.QuotationRequest) error {

	log.Printf("Saving quotation: %+v", req)

	err := rs.dg.SaveQuotationRequest(entity.QuotationEntity{
		RequestId:   req.RequestId,
		UserId:      req.UserId,
		ProductList: req.Products,
		TTL:         strconv.FormatInt(time.Now().Add(time.Minute*5).UnixNano(), 10),
	})
	if err != nil {
		log.Printf("Error saving quotation request: %+v", err)
		return err
	}

	log.Printf("Notifying quotation: %+v", req)
	_, err = rs.snsg.NotifyQuotation(ctx, message.UserPricesMessage{
		RequestId: req.RequestId,
		UserId:    req.UserId,
	})
	if err != nil {
		log.Printf("Error notifying quotation: %+v", err)
		return err
	}

	return nil
}
