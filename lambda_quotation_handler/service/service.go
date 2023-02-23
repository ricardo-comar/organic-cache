package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/ricardo-comar/organic-cache/model"
)

func ResponseWait(ctx *context.Context, cli *dynamodb.Client, requestId string) (*model.QuotationEntity, error) {

	var quotation model.QuotationEntity
	var err error

	// Init hazelcast client
	config := hazelcast.Config{}
	config.Cluster.Network.SetAddresses("hazelcast:5701")
	client, err := hazelcast.StartNewClientWithConfig(*ctx, config)
	if err != nil {
		panic(err)
	}

	// Get a reference to the queue.
	myTopic, err := client.GetTopic(*ctx, "quotation-response-topic")
	if err != nil {
		panic(err)
	}

	// Add a message listener to the topic.
	_, err = myTopic.AddMessageListener(*ctx, func(event *hazelcast.MessagePublished) {
		msg := event.Value.(model.QuotationEntity)
		if msg.Id == requestId {
			quotation = msg
		}
	})
	if err != nil {
		panic(err)
	}

	// Shutdown the client.
	client.Shutdown(*ctx)

	return &quotation, err

}
