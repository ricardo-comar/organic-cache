package gateway

import (
	"context"
	"log"
	"os"

	"github.com/ricardo-comar/organic-cache/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func SaveActiveUser(cli *dynamodb.Client, user model.UserEntity) error {

	log.Printf("Saving user %s", user)
	item, err := attributevalue.MarshalMap(user)

	if err == nil {
		_, err = cli.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("ACTIVE_USERS_TABLE")),
			Item:      item,
		})
	}

	return err
}
