package service

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/model"
)

func GenerateUserPrices(dyncli *dynamodb.Client, user *model.UserEntity) error {

	products, err := gateway.ScanProducts(dyncli)
	if err != nil {
		log.Fatal("Error scanning products :", err)
		return err
	}
	log.Printf("Products: %+v\n", products)

	userDiscounts, err := gateway.QueryUserDiscounts(dyncli, user)
	if err != nil {
		log.Fatal("Error quering user discounts :", err)
		return err
	}
	log.Printf("User discounts: %+v\n", userDiscounts)

	prices := model.UserPricesEntity{
		UserId: user.ID,
	}

	for _, product := range *products {

		finalValue := product.Value
		discount := float32(0.0)

		if userDiscounts != nil {
			for _, userDiscount := range userDiscounts.Discounts {
				if product.ProductId == userDiscount.ProductId {
					discount = userDiscount.Percentage
					finalValue = product.Value * (1 - (discount / 100))
				}
			}
		}

		prices.Prices = append(prices.Prices, model.ProductPrice{
			ProductId:     product.ProductId,
			ProductName:   product.Name,
			OriginalValue: product.Value,
			Value:         finalValue,
			Discount:      discount,
		})
	}

	log.Printf("User prices: %+v\n", prices)

	gateway.SaveUserPrices(dyncli, &prices)
	if err != nil {
		log.Fatal("Error saving user prices :", err)
		return err
	}

	return nil

}
