package service

import (
	"log"

	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/price_calc/gateway"
)

type pricesService struct {
	dg gateway.DynamoGateway
}

func NewPricesService(dg gateway.DynamoGateway) PricesService {
	rs := &pricesService{
		dg: dg,
	}
	return PricesService(rs)
}

type PricesService interface {
	GenerateUserPrices(user *message.UserMessage) error
}

func (ps pricesService) GenerateUserPrices(user *message.UserMessage) error {

	products, err := ps.dg.ScanProducts()
	if err != nil {
		log.Fatal("Error scanning products :", err)
		return err
	}
	log.Printf("Products: %+v\n", products)

	userDiscounts, err := ps.dg.QueryUserDiscounts(user)
	if err != nil {
		log.Fatal("Error quering user discounts :", err)
		return err
	}
	log.Printf("User discounts: %+v\n", userDiscounts)

	prices := entity.UserPricesEntity{
		UserId: user.UserId,
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

		prices.Products = append(prices.Products, entity.ProductPrice{
			ProductId:     product.ProductId,
			ProductName:   product.Name,
			OriginalValue: product.Value,
			Value:         finalValue,
			Discount:      discount,
		})
	}

	log.Printf("User prices: %+v\n", prices)

	ps.dg.SaveUserPrices(&prices)
	if err != nil {
		log.Fatal("Error saving user prices :", err)
		return err
	}

	return nil

}
