package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/ricardo-comar/organic-cache/model"
)

func ResponseWait(ctx *context.Context, cli *dynamodb.Client, requestId string) (*model.QuotationEntity, error) {

	var quotation *model.QuotationEntity
	var err error

	// Init hazelcast client
	config := hazelcast.Config{}
	config.Cluster.Network.SetAddresses("hazelcast:5701")
	client, err := hazelcast.StartNewClientWithConfig(*ctx, config)
	if err != nil {
		return nil, err
	}

	// Get a reference to the queue.
	myTopic, err := client.GetTopic(*ctx, "quotation-response-topic")
	if err != nil {
		return nil, err
	}

	done := make(chan bool, 1)

	// Add a message listener to the topic.
	_, err = myTopic.AddMessageListener(*ctx, func(event *hazelcast.MessagePublished) {
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

	// Shutdown the client.
	client.Shutdown(*ctx)

	return quotation, err

}
