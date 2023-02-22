package gateway

import (
	"context"
	"os"

	"github.com/ricardo-comar/organic-cache/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func QueryQuotation(cli *dynamodb.Client, requestId string) (*model.ProductQuotation, error) {

	output, err := cli.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("QUOTATION_TABLE")),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: requestId},
		},
	})

	var quotation *model.ProductQuotation
	if err == nil && output.Item != nil {
		attributevalue.UnmarshalMap(output.Item, &quotation)
	}

	return quotation, nil
}
