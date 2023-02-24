package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/model"
)

func ResponseWait(ctx *context.Context, requestId string) (*model.QuotationEntity, error) {

	var quotation *model.QuotationEntity
	var err error
	done := make(chan bool, 1)

	err = gateway.SubscribeTopic(ctx, "quotation-response-topic", func(message []byte) {
		msg := model.QuotationEntity{}
		json.Unmarshal(message, &msg)

		if msg.Id == requestId {
			quotation = &msg
			done <- true
		}
	})

	// Add a timeout
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()

	<-done

	return quotation, err

}
