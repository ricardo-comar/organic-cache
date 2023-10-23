package gateway

import (
	"context"
	"os"

	"github.com/ricardo-comar/organic-cache/cache_refresh/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func QueryUsers(cli *dynamodb.Client) ([]model.UserEntity, error) {

	var totalUsers []model.UserEntity

	input := dynamodb.NewScanPaginator(cli, &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("ACTIVE_USERS_TABLE")),
	})

	for input.HasMorePages() {
		out, err := input.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		var users []model.UserEntity
		err = attributevalue.UnmarshalListOfMaps(out.Items, &users)
		if err != nil {
			return nil, err
		}

		totalUsers = append(totalUsers, users...)
	}

	return totalUsers, nil
}
