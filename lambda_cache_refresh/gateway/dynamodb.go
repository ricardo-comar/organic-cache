package gateway

import (
	"context"
	"os"

	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/gateway"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	QueryUsers() ([]entity.UserEntity, error)
}

func (dg dynamoCxt) QueryUsers() ([]entity.UserEntity, error) {

	var totalUsers []entity.UserEntity

	input := dynamodb.NewScanPaginator(dg.dyncli, &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("ACTIVE_USERS_TABLE")),
	})

	for input.HasMorePages() {
		out, err := input.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		var users []entity.UserEntity
		err = attributevalue.UnmarshalListOfMaps(out.Items, &users)
		if err != nil {
			return nil, err
		}

		totalUsers = append(totalUsers, users...)
	}

	return totalUsers, nil
}
