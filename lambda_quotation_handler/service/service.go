package service

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/model"
)

func ResponseWait(ctx *context.Context, cli *dynamodb.Client, requestId string) (*model.ProductQuotation, error) {

	var quotation *model.ProductQuotation
	var err error

	for i := 0; i < 10; i++ {
		quotation, err = gateway.QueryQuotation(cli, requestId)

		if quotation == nil {
			time.Sleep(time.Second)
		}
	}

	return quotation, err

}
