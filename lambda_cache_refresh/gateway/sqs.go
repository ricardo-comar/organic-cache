package gateway

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SendMessage(ctx context.Context, cfg aws.Config, message interface{}) (*string, error) {

	svc := sqs.NewFromConfig(cfg)
	localendpoint, found := os.LookupEnv("LOCALSTACK_HOSTNAME")
	if found {
		svc = sqs.New(sqs.Options{Credentials: cfg.Credentials, EndpointResolver: sqs.EndpointResolverFromURL("http://" + localendpoint + ":" + os.Getenv("EDGE_PORT"))})
	}

	body, _ := json.Marshal(message)
	res, err := svc.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(os.Getenv("REFRESH_QUEUE")),
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res.MessageId, nil
}
