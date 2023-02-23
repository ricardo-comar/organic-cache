package service

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardo-comar/identity-provider/gateway"
	"github.com/ricardo-comar/identity-provider/model"
)

func GenerateUserPrices(dyncli *dynamodb.Client, user *model.UserEntity) error {

	products, err := gateway.ScanProducts(dyncli)
	if err != nil {
		log.Fatal("Error scanning products :", err)
		return err
	}
	fmt.Printf("Products: %+v\n", products)

	userDiscounts, err := gateway.QueryUserDiscounts(dyncli, user)
	if err != nil {
		log.Fatal("Error quering user discounts :", err)
		return err
	}
	fmt.Printf("User discounts: %+v\n", userDiscounts)

	prices := model.UserPricesEntity{}
	prices.UserId = user.ID
	prices.Prices = []model.ProductPrice{}

	for _, product := range *products {

		finalValue := product.Value
		discount := float32(0.0)

		for _, userDiscount := range userDiscounts.Discounts {
			if product.ProductId == userDiscount.ProductId {
				discount = userDiscount.Percentage
				finalValue = product.Value * (discount / 100)
			}
		}

		prices.Prices = append(prices.Prices, model.ProductPrice{
			ProductId:     product.ProductId,
			OriginalValue: product.Value,
			Value:         finalValue,
			Discount:      discount,
		})
	}

	fmt.Printf("User prices: %+v\n", prices)

	gateway.SaveUserPrices(dyncli, &prices)
	if err != nil {
		log.Fatal("Error saving user prices :", err)
		return err
	}

	return nil

}
