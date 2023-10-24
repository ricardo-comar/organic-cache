package gateway

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
)

func NotifyQuotation(ctx context.Context, cli *sns.Client, msg message.UserPricesMessage) (*string, error) {

	body, _ := json.Marshal(msg)
	res, err := cli.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(string(body)),
		TopicArn: aws.String(os.Getenv("USER_PRICES_TOPIC_ARN")),
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res.MessageId, nil
}
