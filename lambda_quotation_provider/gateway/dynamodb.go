package gateway

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ricardo-comar/identity-provider/model"
)

func QueryProductPrice(cli dynamodb.Client, userId string) ([]model.ProductPrice, error) {

	items, err := cli.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("USER_PRICES_TABLE")),
		KeyConditionExpression: aws.String("user_id = :userKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userKey": &types.AttributeValueMemberS{Value: userId},
		},
	})

	if err != nil {
		panic(err)
	}

	var products []model.ProductPrice
	err = attributevalue.UnmarshalListOfMaps(items.Items, &products)
	if err != nil {
		panic(err)
	}

	return products, nil

}
