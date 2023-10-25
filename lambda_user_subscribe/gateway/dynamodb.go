package gateway

import (
	"context"
	"log"
	"os"

	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/gateway"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	SaveActiveUser(user *entity.UserEntity) error
	QuerySubscription(userId string) (*entity.UserEntity, error)
}

func (dg dynamoCxt) SaveActiveUser(user *entity.UserEntity) error {

	log.Printf("Saving user %s", user)
	item, err := attributevalue.MarshalMap(user)

	if err == nil {
		_, err = dg.dyncli.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("ACTIVE_USERS_TABLE")),
			Item:      item,
		})
	}

	return err
}

func (dg dynamoCxt) QuerySubscription(userId string) (*entity.UserEntity, error) {

	output, err := dg.dyncli.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("ACTIVE_USERS_TABLE")),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: userId},
		},
	})

	if err == nil && output.Item != nil {
		userSub := entity.UserEntity{}
		err = attributevalue.UnmarshalMap(output.Item, &userSub)
		return &userSub, err
	}

	return nil, err

}
