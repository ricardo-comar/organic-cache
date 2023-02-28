package gateway

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/ricardo-comar/organic-cache/model"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

func SendResponse(gtwcli *apigatewaymanagementapi.Client, ctx context.Context, quotation model.QuotationEntity, connectionId string) (*apigatewaymanagementapi.PostToConnectionOutput, error) {

	data, err := json.Marshal(quotation)

	resp, err := gtwcli.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
		Data:         data,
		ConnectionId: aws.String(connectionId),
	})

	return resp, err
}
