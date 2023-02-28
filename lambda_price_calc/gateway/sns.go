package gateway

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func NotifyQuotation(ctx context.Context, cli *sns.Client, requestId *string) (*string, error) {

	res, err := cli.Publish(ctx, &sns.PublishInput{
		Message:  aws.String("{ \"requestId\": \"" + *requestId + "\"}"),
		TopicArn: aws.String(os.Getenv("QUOTATIONS_TOPIC_ARN")),
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res.MessageId, nil
}
