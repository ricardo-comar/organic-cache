package gateway

import (
	"context"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/ricardo-comar/identity-provider/model"
)

func SendResponse(ctx *context.Context, quotation model.QuotationEntity) {
	// Init hazelcast client
	config := hazelcast.Config{}
	config.Cluster.Network.SetAddresses("localhost:5701")
	client, err := hazelcast.StartNewClientWithConfig(*ctx, config)
	if err != nil {
		panic(err)
	}

	// Get a reference to the queue.
	myTopic, err := client.GetTopic(*ctx, "quotation-response-topic")
	if err != nil {
		panic(err)
	}

	err = myTopic.Publish(*ctx, quotation)
	if err != nil {
		panic(err)
	}

}
