package gateway

import (
	"context"

	"github.com/hazelcast/hazelcast-go-client"
)

func HazelcastTopic(ctx *context.Context) (*hazelcast.Topic, error) {

	// Init hazelcast client
	config := hazelcast.Config{}
	config.Cluster.Network.SetAddresses("hazelcast:5701")
	client, err := hazelcast.StartNewClientWithConfig(*ctx, config)
	if err != nil {
		return nil, err
	}

	// Get a reference to the queue.
	topic, err := client.GetTopic(*ctx, "quotation-response-topic")
	if err != nil {
		return nil, err
	}

	return topic, nil
}
