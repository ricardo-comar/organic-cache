package gateway

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ricardo-comar/organic-cache/model"
)

func QueryProductPrice(cli dynamodb.Client, userId string) (*model.UserPricesEntity, error) {

	output, err := cli.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("USER_PRICES_TABLE")),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userId},
		},
	})

	if err == nil && output.Item != nil {
		userPrices := model.UserPricesEntity{}
		err = attributevalue.UnmarshalMap(output.Item, &userPrices)
		return &userPrices, err
	}

	return nil, err

}

func QueryRequest(cli dynamodb.Client, requestId string) (*model.QuotationRequest, error) {

	output, err := cli.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("QUOTATIONS_TABLE")),
		Key: map[string]types.AttributeValue{
			"request_id": &types.AttributeValueMemberS{Value: requestId},
		},
	})

	if err == nil && output.Item != nil {
		request := model.QuotationRequest{}
		err = attributevalue.UnmarshalMap(output.Item, &request)
		return &request, err
	}

	return nil, err

}
