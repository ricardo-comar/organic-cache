package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardo-comar/organic-cache/lib_common/api"
	"github.com/ricardo-comar/organic-cache/quotation_handler/gateway"
	"github.com/ricardo-comar/organic-cache/quotation_handler/service"
)

func main() {
	lambdaHandler := awsHandler{
		reqSrv: service.NewRequestService(gateway.NewDynamoGateway(), gateway.NewSNSGateway()),
		rspSrv: service.NewResponseService(gateway.NewSQSGateway()),
	}
	lambda.Start(lambdaHandler.handleRequest)
}

type awsHandler struct {
	reqSrv service.RequestService
	rspSrv service.ResponseService
}

func (l awsHandler) handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	req := api.QuotationRequest{RequestId: request.RequestContext.RequestID}
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error parsing request: %+v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}
	log.Printf("Message: %+v", req)

	err = l.reqSrv.RequestQuotation(ctx, req)
	if err != nil {
		log.Printf("Error saving quotation request: %+v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	response, err := l.rspSrv.WaitForResponse(ctx, req.RequestId)

	if err == nil {
		log.Printf("Response: %+v", response)
		resp, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{Body: string(resp), StatusCode: http.StatusOK}, nil
	}

	switch err.(type) {
	case *service.TimeoutError:
		log.Printf("Timeout waiting for quotation response: %+v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusRequestTimeout}, err
	}

	log.Printf("Error waiting for quotation response: %+v", err)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
}
