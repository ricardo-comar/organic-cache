package gateway

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

func SendResponse(ctx context.Context, cli *sns.Client, msgId *string, message interface{}) (*string, error) {

	body, _ := json.Marshal(message)
	res, err := cli.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(string(body)),
		TopicArn: aws.String(os.Getenv("RESPONSE_TOPIC_ARN")),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"": {
				DataType:    aws.String("String"),
				StringValue: msgId,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res.MessageId, nil
}
