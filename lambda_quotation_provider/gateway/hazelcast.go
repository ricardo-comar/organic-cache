package gateway

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/ricardo-comar/organic-cache/model"
)

func SendResponse(ctx *context.Context, quotation *model.QuotationEntity) {
	// Init hazelcast client
	config := hazelcast.Config{}
	config.Cluster.Network.SetAddresses("hazelcast:5701")
	client, err := hazelcast.StartNewClientWithConfig(*ctx, config)
	if err != nil {
		log.Printf("Erro ao conectar no hazelcast: %+v", err)
		return
	}

	// Get a reference to the queue.
	myTopic, err := client.GetTopic(*ctx, "quotation-response-topic")
	if err != nil {
		log.Printf("Erro ao conectar no tópico: %+v", err)
		return
	}

	msg, _ := json.Marshal(*quotation)
	err = myTopic.Publish(*ctx, msg)
	if err != nil {
		log.Printf("Erro ao enviar mensagem: %+v", err)
		return
	}

}
