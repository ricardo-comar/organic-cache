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

func QueryUserDiscounts(dyncli *dynamodb.Client, user *model.UserEntity) (*model.DiscountEntity, error) {

	output, err := dyncli.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("USER_DISCOUNTS_TABLE")),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: user.ID},
		},
	})

	discounts := model.DiscountEntity{}
	if err == nil && output.Item != nil {
		attributevalue.UnmarshalMap(output.Item, &discounts)
	}

	return &discounts, nil

}

func ScanProducts(dyncli *dynamodb.Client) (*[]model.ProductEntity, error) {

	var totalProducts []model.ProductEntity

	input := dynamodb.NewScanPaginator(dyncli, &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("PRODUCTS_TABLE")),
	})

	for input.HasMorePages() {
		out, err := input.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}

		var products []model.ProductEntity
		err = attributevalue.UnmarshalListOfMaps(out.Items, &products)
		if err != nil {
			panic(err)
		}

		totalProducts = append(totalProducts, products...)
	}

	return &totalProducts, nil
}

func SaveUserPrices(dyncli *dynamodb.Client, prices *model.UserPricesEntity) error {

	item, err := attributevalue.MarshalMap(prices)

	if err == nil {
		_, err = dyncli.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("USER_PRICES_TABLE")),
			Item:      item,
		})
	}

	return err

}
