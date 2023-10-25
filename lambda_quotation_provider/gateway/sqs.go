package gateway

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
)

func RecalcMessage(ctx context.Context, cli *sqs.Client, msg *message.UserPricesMessage) (*string, error) {

	body, _ := json.Marshal(msg)
	res, err := cli.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:    aws.String(string(body)),
		QueueUrl:       aws.String(os.Getenv("RECALC_QUEUE")),
		MessageGroupId: aws.String("quotation"),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"RequestId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(msg.RequestId),
			},
		},
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res.MessageId, nil
}

func NotifyQuotation(ctx context.Context, cli *sqs.Client, msg message.QuotationMessage) (*string, error) {

	body, _ := json.Marshal(msg)
	res, err := cli.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(os.Getenv("QUOTATIONS_QUEUE")),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"RequestId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(msg.RequestId),
			},
		},
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res.MessageId, nil
}
