package gateway

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
)

func SaveQuotationRequest(dyncli *dynamodb.Client, req entity.QuotationEntity) error {

	item, err := attributevalue.MarshalMap(req)

	if err == nil {
		_, err = dyncli.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("QUOTATIONS_TABLE")),
			Item:      item,
		})
	}

	return err

}
