package gateway

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func ReceiveMessage(ctx context.Context, sqscli *sqs.Client) (*sqs.ReceiveMessageOutput, error) {

	result, err := sqscli.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(os.Getenv("QUOTATIONS_QUEUE")),
		MaxNumberOfMessages:   1,
		WaitTimeSeconds:       3,
		AttributeNames:        []types.QueueAttributeName{"All"},
		MessageAttributeNames: []string{"All"},
	})

	return result, err
}

func ChangeMessageVisibility(ctx context.Context, sqscli *sqs.Client, msgReceiptHandle *string) error {

	_, err := sqscli.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(os.Getenv("QUOTATIONS_QUEUE")),
		ReceiptHandle:     msgReceiptHandle,
		VisibilityTimeout: 1,
	})

	return err
}

func DeleteMessage(ctx context.Context, sqscli *sqs.Client, msgReceiptHandle *string) error {

	_, err := sqscli.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(os.Getenv("QUOTATIONS_QUEUE")),
		ReceiptHandle: msgReceiptHandle,
	})

	return err
}
