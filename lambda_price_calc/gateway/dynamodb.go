package gateway

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardo-comar/identity-provider/model"
)

func SaveActiveUser(svc dynamodb.Client, user model.UserEntity) error {

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
