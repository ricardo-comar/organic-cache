package service_test

import (
	"errors"
	"testing"

	"github.com/go-test/deep"
	"github.com/ricardo-comar/organic-cache/lib_common/entity"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/price_calc/service"
	"github.com/stretchr/testify/assert"
)

var userMsg *message.UserMessage
var dynProducts *[]entity.ProductEntity
var dynUserDiscounts *entity.DiscountEntity
var userPrices entity.UserPricesEntity
var userPricesNoDisc entity.UserPricesEntity
var calledQueryUserDiscounts int
var calledScanProducts int
var calledSaveUserPrices int

func initVariables() {
	calledQueryUserDiscounts = 0
	calledScanProducts = 0
	calledSaveUserPrices = 0

	userMsg = &message.UserMessage{
		UserId: "ABC",
	}

	dynProducts = &[]entity.ProductEntity{
		{ProductId: "P1", Name: "Mock Product 1", Value: 134.22},
		{ProductId: "P2", Name: "Mock Product 2", Value: 6500},
		{ProductId: "P3", Name: "Mock Product 3", Value: 987.5},
	}
	dynUserDiscounts = &entity.DiscountEntity{
		UserId: userMsg.UserId,
		Discounts: []entity.ProductDiscount{
			{ProductId: "P1", Percentage: 15},
			{ProductId: "P3", Percentage: 7},
		},
	}
	userPrices = entity.UserPricesEntity{
		UserId: userMsg.UserId,
		Products: []entity.ProductPrice{
			{
				ProductId:     (*dynProducts)[0].ProductId,
				ProductName:   (*dynProducts)[0].Name,
				OriginalValue: (*dynProducts)[0].Value,
				Value:         (*dynProducts)[0].Value * (1 - ((*dynUserDiscounts).Discounts[0].Percentage / 100)),
				Discount:      (*dynUserDiscounts).Discounts[0].Percentage,
			},
			{
				ProductId:     (*dynProducts)[1].ProductId,
				ProductName:   (*dynProducts)[1].Name,
				OriginalValue: (*dynProducts)[1].Value,
				Value:         (*dynProducts)[1].Value,
				Discount:      0,
			},
			{
				ProductId:     (*dynProducts)[2].ProductId,
				ProductName:   (*dynProducts)[2].Name,
				OriginalValue: (*dynProducts)[2].Value,
				Value:         (*dynProducts)[2].Value * (1 - ((*dynUserDiscounts).Discounts[1].Percentage / 100)),
				Discount:      (*dynUserDiscounts).Discounts[1].Percentage,
			},
		},
	}
	userPricesNoDisc = entity.UserPricesEntity{
		UserId: userMsg.UserId,
		Products: []entity.ProductPrice{
			{
				ProductId:     (*dynProducts)[0].ProductId,
				ProductName:   (*dynProducts)[0].Name,
				OriginalValue: (*dynProducts)[0].Value,
				Value:         (*dynProducts)[0].Value,
				Discount:      0,
			},
			{
				ProductId:     (*dynProducts)[1].ProductId,
				ProductName:   (*dynProducts)[1].Name,
				OriginalValue: (*dynProducts)[1].Value,
				Value:         (*dynProducts)[1].Value,
				Discount:      0,
			},
			{
				ProductId:     (*dynProducts)[2].ProductId,
				ProductName:   (*dynProducts)[2].Name,
				OriginalValue: (*dynProducts)[2].Value,
				Value:         (*dynProducts)[2].Value,
				Discount:      0,
			},
		},
	}

}
func TestSuccessRequest(t *testing.T) {
	initVariables()

	tstService := service.NewPricesService(mockDynamoGateway{t: t, dynProducts: dynProducts, discount: dynUserDiscounts, user: userMsg, userPrices: &userPrices})

	respErr := tstService.GenerateUserPrices(userMsg)
	assert.Nil(t, respErr, "Unexpected Error")

	assert.Equal(t, 1, calledQueryUserDiscounts, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 1, calledScanProducts, "Unexpected ScanProducts calls")
	assert.Equal(t, 1, calledSaveUserPrices, "Unexpected SaveUserPrices calls")

}
func TestErrorScanProducts(t *testing.T) {
	initVariables()

	errScanProducts := errors.New("errScanProducts")

	tstService := service.NewPricesService(mockDynamoGateway{t: t, errScanProducts: errScanProducts})

	respErr := tstService.GenerateUserPrices(userMsg)
	assert.Equal(t, errScanProducts, respErr, "Unexpected Error")

	assert.Equal(t, 0, calledQueryUserDiscounts, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 1, calledScanProducts, "Unexpected ScanProducts calls")
	assert.Equal(t, 0, calledSaveUserPrices, "Unexpected SaveUserPrices calls")

}
func TestNilScanProducts(t *testing.T) {
	initVariables()

	tstService := service.NewPricesService(mockDynamoGateway{t: t, dynProducts: nil})

	respErr := tstService.GenerateUserPrices(userMsg)
	assert.NotNil(t, respErr, "Error is nil")

	assert.Equal(t, 0, calledQueryUserDiscounts, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 1, calledScanProducts, "Unexpected ScanProducts calls")
	assert.Equal(t, 0, calledSaveUserPrices, "Unexpected SaveUserPrices calls")

}
func TestEmptyScanProducts(t *testing.T) {
	initVariables()

	tstService := service.NewPricesService(mockDynamoGateway{t: t, dynProducts: &[]entity.ProductEntity{}})

	respErr := tstService.GenerateUserPrices(userMsg)
	assert.NotNil(t, respErr, "Error is nil")

	assert.Equal(t, 0, calledQueryUserDiscounts, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 1, calledScanProducts, "Unexpected ScanProducts calls")
	assert.Equal(t, 0, calledSaveUserPrices, "Unexpected SaveUserPrices calls")

}
func TestNilQueryUserDiscounts(t *testing.T) {
	initVariables()

	tstService := service.NewPricesService(mockDynamoGateway{t: t, dynProducts: dynProducts, discount: nil, user: userMsg})

	respErr := tstService.GenerateUserPrices(userMsg)
	assert.NotNil(t, respErr, "Error is nil")

	assert.Equal(t, 1, calledQueryUserDiscounts, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 1, calledScanProducts, "Unexpected ScanProducts calls")
	assert.Equal(t, 0, calledSaveUserPrices, "Unexpected SaveUserPrices calls")

}
func TestErrorQueryUserDiscounts(t *testing.T) {
	initVariables()
	errQueryUserDiscounts := errors.New("errQueryUserDiscounts")

	tstService := service.NewPricesService(mockDynamoGateway{t: t, dynProducts: dynProducts, user: userMsg, errQueryUserDiscounts: errQueryUserDiscounts})

	respErr := tstService.GenerateUserPrices(userMsg)
	assert.Equal(t, respErr, errQueryUserDiscounts, "Error is unpexpected")

	assert.Equal(t, 1, calledQueryUserDiscounts, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 1, calledScanProducts, "Unexpected ScanProducts calls")
	assert.Equal(t, 0, calledSaveUserPrices, "Unexpected SaveUserPrices calls")

}
func TestEmptyQueryUserDiscounts(t *testing.T) {
	initVariables()

	tstService := service.NewPricesService(mockDynamoGateway{t: t, dynProducts: dynProducts, discount: &entity.DiscountEntity{}, user: userMsg, userPrices: &userPricesNoDisc})

	respErr := tstService.GenerateUserPrices(userMsg)
	assert.Nil(t, respErr, "Unexpected Error")

	assert.Equal(t, 1, calledQueryUserDiscounts, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 1, calledScanProducts, "Unexpected ScanProducts calls")
	assert.Equal(t, 1, calledSaveUserPrices, "Unexpected SaveUserPrices calls")

}

func TestErrorSaveUserPrices(t *testing.T) {
	initVariables()
	errSaveUserPrices := errors.New("errSaveUserPrices")

	tstService := service.NewPricesService(mockDynamoGateway{t: t, dynProducts: dynProducts, user: userMsg, errSaveUserPrices: errSaveUserPrices, discount: dynUserDiscounts, userPrices: &userPrices})

	respErr := tstService.GenerateUserPrices(userMsg)
	assert.Equal(t, errSaveUserPrices, respErr, "Error is unpexpected")

	assert.Equal(t, 1, calledQueryUserDiscounts, "Unexpected QueryUserDiscounts calls")
	assert.Equal(t, 1, calledScanProducts, "Unexpected ScanProducts calls")
	assert.Equal(t, 1, calledSaveUserPrices, "Unexpected SaveUserPrices calls")

}

type mockDynamoGateway struct {
	t                     *testing.T
	dynProducts           *[]entity.ProductEntity
	user                  *message.UserMessage
	discount              *entity.DiscountEntity
	userPrices            *entity.UserPricesEntity
	errQueryUserDiscounts error
	errScanProducts       error
	errSaveUserPrices     error
}

func (g mockDynamoGateway) QueryUserDiscounts(user *message.UserMessage) (*entity.DiscountEntity, error) {
	calledQueryUserDiscounts++

	if g.user == nil {
		assert.Nil(g.t, user, "Unexpected user")
	} else {
		if diff := deep.Equal(g.user, user); diff != nil {
			g.t.Error(diff)
		}
	}

	return g.discount, g.errQueryUserDiscounts
}
func (g mockDynamoGateway) ScanProducts() (*[]entity.ProductEntity, error) {
	calledScanProducts++

	return g.dynProducts, g.errScanProducts
}
func (g mockDynamoGateway) SaveUserPrices(prices *entity.UserPricesEntity) error {
	calledSaveUserPrices++

	if g.userPrices == nil {
		assert.Nil(g.t, prices, "Unexpected user")
	} else {
		if diff := deep.Equal(g.userPrices, prices); diff != nil {
			g.t.Error(diff)
		}
	}

	return g.errSaveUserPrices
}
