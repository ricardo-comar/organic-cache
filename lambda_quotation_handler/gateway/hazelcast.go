package gateway

import (
	"context"

	"github.com/hazelcast/hazelcast-go-client"
)

type MessageHandler func(message []byte)

func SubscribeTopic(ctx *context.Context, topicName string, callback MessageHandler) error {

	// Init hazelcast client
	config := hazelcast.Config{}
	config.Cluster.Network.SetAddresses("hazelcast:5701")
	client, err := hazelcast.StartNewClientWithConfig(*ctx, config)
	if err != nil {
		return err
	}

	// Get a reference to the queue.
	topic, err := client.GetTopic(*ctx, topicName)

	_, err = topic.AddMessageListener(*ctx, func(event *hazelcast.MessagePublished) {
		callback(event.Value.([]byte)[:])
	})

	return err
}
