package gateway

import (
	"context"
	"fmt"
	"os"

	"github.com/ricardo-comar/organic-cache/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func SaveActiveUser(cfg aws.Config, user model.UserEntity) error {

	svc := dynamodb.NewFromConfig(cfg)

	localendpoint, found := os.LookupEnv("LOCALSTACK_HOSTNAME")
	if found {
		svc = dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL("http://"+localendpoint+":4566")))
	}

	fmt.Printf("Saving user %s", user)
	item, err := attributevalue.MarshalMap(user)

	if err == nil {
		_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("ACTIVE_USERS_TABLE")),
			Item:      item,
		})
	}

	return err
}
