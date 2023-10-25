package gateway

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/gateway"
)

type dynamoCxt struct {
	dyncli *dynamodb.Client
}

func NewDynamoGateway() DynamoGateway {
	dg := &dynamoCxt{dyncli: gateway.InitDynamoClient()}
	gtw := DynamoGateway(dg)

	return gtw
}

type DynamoGateway interface {
	SaveQuotationRequest(req entity.QuotationEntity) error
}

func (dg dynamoCxt) SaveQuotationRequest(req entity.QuotationEntity) error {

	item, err := attributevalue.MarshalMap(req)

	if err == nil {
		_, err = dg.dyncli.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("QUOTATIONS_TABLE")),
			Item:      item,
		})
	}

	return err

}
