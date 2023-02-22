package service

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func ResponseWait(ctx context.Context, cli *sns.Client, msgId *string) (string, error) {

	resultSub, _ := cli.Subscribe(ctx, &sns.SubscribeInput{
		TopicArn: aws.String(os.Getenv("RESPONSE_TOPIC_ARN")),
		Attributes: map[string]string{
			"FilterPolicy": "{ \"RequestId\": \"" + *msgId + "\"}",
		},
	})

	fmt.Printf("resultSub: %+v\n", resultSub)

	return "nil", nil
}
