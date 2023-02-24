package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/ricardo-comar/organic-cache/gateway"
	"github.com/ricardo-comar/organic-cache/model"
)

func ResponseWait(ctx *context.Context, requestId string) (*model.QuotationEntity, error) {

	var quotation *model.QuotationEntity
	var err error
	topic, err := gateway.HazelcastTopic(ctx)

	done := make(chan bool, 1)

	// Add a message listener to the topic.
	_, err = topic.AddMessageListener(*ctx, func(event *hazelcast.MessagePublished) {
		log.Printf("**** Message received: %+v", string(event.Value.([]byte)[:]))

		msg := model.QuotationEntity{}
		json.Unmarshal(event.Value.([]byte), &msg)

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
