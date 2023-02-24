package gateway

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func SubscribeUser(ctx context.Context, cli *sns.Client, message interface{}) (*string, error) {

	body, _ := json.Marshal(message)
	res, err := cli.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(string(body)),
		TopicArn: aws.String(os.Getenv("USER_SUBSCRIBE_TOPIC")),
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res.MessageId, nil
}

type UserMessage struct {
	ID string `json:"id"`
}
