package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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
	RecalcMessage(ctx context.Context, msg *message.UserPricesMessage) (*string, error)
	NotifyQuotation(ctx context.Context, msg *message.QuotationMessage) (*string, error)
}

func (gtw sqsCxt) RecalcMessage(ctx context.Context, msg *message.UserPricesMessage) (*string, error) {

	body, _ := json.Marshal(msg)
	res, err := gtw.sqscli.SendMessage(ctx, &sqs.SendMessageInput{
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

func (gtw sqsCxt) NotifyQuotation(ctx context.Context, msg *message.QuotationMessage) (*string, error) {

	if msg == nil {
		log.Println("QuotationMessage is nil")
		return nil, errors.New("nil_quotation_message")
	}

	body, _ := json.Marshal(msg)
	res, err := gtw.sqscli.SendMessage(ctx, &sqs.SendMessageInput{
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
