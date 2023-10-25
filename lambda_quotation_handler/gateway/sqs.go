package gateway

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ricardo-comar/organic-cache/lib_common/gateway"
)

type sqsCxt struct {
	sqscli *sqs.Client
}

func NewSQSGateway() SQSGateway {
	ctx := &sqsCxt{sqscli: gateway.InitSQSClient()}
	gtw := SQSGateway(ctx)

	return gtw
}

type SQSGateway interface {
	ReceiveMessage(ctx context.Context) (*sqs.ReceiveMessageOutput, error)
	ChangeMessageVisibility(ctx context.Context, msgReceiptHandle *string) error
	DeleteMessage(ctx context.Context, msgReceiptHandle *string) error
}

func (gtw sqsCxt) ReceiveMessage(ctx context.Context) (*sqs.ReceiveMessageOutput, error) {

	result, err := gtw.sqscli.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(os.Getenv("QUOTATIONS_QUEUE")),
		MaxNumberOfMessages:   1,
		WaitTimeSeconds:       3,
		AttributeNames:        []types.QueueAttributeName{"All"},
		MessageAttributeNames: []string{"All"},
	})

	return result, err
}

func (gtw sqsCxt) ChangeMessageVisibility(ctx context.Context, msgReceiptHandle *string) error {

	_, err := gtw.sqscli.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(os.Getenv("QUOTATIONS_QUEUE")),
		ReceiptHandle:     msgReceiptHandle,
		VisibilityTimeout: 1,
	})

	return err
}

func (gtw sqsCxt) DeleteMessage(ctx context.Context, msgReceiptHandle *string) error {

	_, err := gtw.sqscli.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(os.Getenv("QUOTATIONS_QUEUE")),
		ReceiptHandle: msgReceiptHandle,
	})

	return err
}
