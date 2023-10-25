package gateway

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ricardo-comar/organic-cache/lib_common/gateway"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
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
	SendMessage(ctx context.Context, user *message.UserMessage) (*string, error)
}

func (gtw sqsCxt) SendMessage(ctx context.Context, user *message.UserMessage) (*string, error) {

	body, _ := json.Marshal(user)
	res, err := gtw.sqscli.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:    aws.String(string(body)),
		QueueUrl:       aws.String(os.Getenv("RECALC_QUEUE")),
		MessageGroupId: aws.String("price-refresh"),
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res.MessageId, nil
}
