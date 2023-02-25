package gateway

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func RetryMessage(ctx context.Context, cli *sqs.Client, msgId *string, retryCount *string, message interface{}) (*string, error) {

	body, _ := json.Marshal(message)
	res, err := cli.SendMessage(ctx, &sqs.SendMessageInput{
		MessageAttributes: map[string]types.MessageAttributeValue{
			"RequestId": {
				DataType:    aws.String("String"),
				StringValue: msgId,
			},
			"RetryCount": {
				DataType:    aws.String("String"),
				StringValue: retryCount,
			},
		},
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(os.Getenv("QUOTATION_QUEUE")),
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res.MessageId, nil
}
func RecalcMessage(ctx context.Context, cli *sqs.Client, userId string) (*string, error) {

	payload, err := json.Marshal(map[string]interface{}{
		"id": userId,
	})

	res, err := cli.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:    aws.String(string(payload)),
		QueueUrl:       aws.String(os.Getenv("RECALC_QUEUE")),
		MessageGroupId: aws.String("quotation"),
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res.MessageId, nil
}
