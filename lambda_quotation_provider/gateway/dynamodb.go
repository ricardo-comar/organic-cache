package gateway

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	QueryProductPrice(userId string) (*entity.UserPricesEntity, error)
	QueryRequest(requestId string) (*entity.QuotationEntity, error)
}

func (dg dynamoCxt) QueryProductPrice(userId string) (*entity.UserPricesEntity, error) {

	output, err := dg.dyncli.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("USER_PRICES_TABLE")),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userId},
		},
	})

	if err == nil && output.Item != nil {
		userPrices := entity.UserPricesEntity{}
		err = attributevalue.UnmarshalMap(output.Item, &userPrices)
		return &userPrices, err
	}

	return nil, err

}

func (dg dynamoCxt) QueryRequest(requestId string) (*entity.QuotationEntity, error) {

	output, err := dg.dyncli.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("QUOTATIONS_TABLE")),
		Key: map[string]types.AttributeValue{
			"request_id": &types.AttributeValueMemberS{Value: requestId},
		},
	})

	if err == nil && output.Item != nil {
		request := entity.QuotationEntity{}
		err = attributevalue.UnmarshalMap(output.Item, &request)
		return &request, err
	}

	return nil, err

}
