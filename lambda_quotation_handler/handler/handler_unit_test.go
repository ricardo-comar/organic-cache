package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/ricardo-comar/organic-cache/lib_common/api"
	"github.com/ricardo-comar/organic-cache/lib_common/model"
	"github.com/ricardo-comar/organic-cache/quotation_handler/service"
	"github.com/stretchr/testify/assert"
)

var userId string
var productQuotation model.ProductQuotation
var gatewayRequestContext events.APIGatewayProxyRequestContext
var apiRequest api.QuotationRequest
var gatewayRequest events.APIGatewayProxyRequest
var apiResponse api.QuotationResponse
var gatewayResponse events.APIGatewayProxyResponse
var ctx context.Context

func initVariables() {

	userId = "MockUser"
	productQuotation = model.ProductQuotation{
		ProductId:     "A",
		Quantity:      1,
		ProductName:   "Mock Product A",
		OriginalValue: 100.0,
		Discount:      10.0,
		FinalValue:    90.0,
	}

	gatewayRequestContext = events.APIGatewayProxyRequestContext{
		RequestID: uuid.New().String(),
	}

	apiRequest = api.QuotationRequest{
		RequestId: gatewayRequestContext.RequestID,
		UserId:    userId,
		Products: []model.QuotationItem{
			{
				ProductId: productQuotation.ProductId,
				Quantity:  productQuotation.Quantity,
			},
		},
	}
	body, _ := json.Marshal(apiRequest)

	gatewayRequest = events.APIGatewayProxyRequest{
		RequestContext: gatewayRequestContext,
		Body:           string(body),
	}

	apiResponse = api.QuotationResponse{
		RequestId: gatewayRequestContext.RequestID,
		UserId:    userId,
		Products: []model.ProductQuotation{
			productQuotation,
		},
	}
	respBody, _ := json.Marshal(apiResponse)

	ctx = context.TODO()

	gatewayResponse = events.APIGatewayProxyResponse{Body: string(respBody), StatusCode: http.StatusOK}

}
func TestSuccess(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		reqSrv: mockRequestService{t: t, ctx: ctx, apiReq: apiRequest},
		rspSrv: mockResponseService{t: t, ctx: ctx, requestId: gatewayRequestContext.RequestID, resp: &apiResponse},
	}

	response, err := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.Nil(t, err, "Unexpected error")
	assert.NotNil(t, response, "Nil response")
	if diff := deep.Equal(response, gatewayResponse); diff != nil {
		t.Error(diff)
	}

}

func TestInvalidRequestBody(t *testing.T) {
	initVariables()

	lambdaHandler := awsHandler{
		reqSrv: mockRequestService{t: t, ctx: ctx, apiReq: apiRequest},
		rspSrv: mockResponseService{t: t, ctx: ctx, requestId: gatewayRequestContext.RequestID, resp: &apiResponse},
	}

	gatewayRequest.Body = "ADI&%ASDB"
	gatewayResponse = events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}

	response, err := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.NotNil(t, err, "Unexpected error")
	assert.NotNil(t, response, "Nil response")
	if diff := deep.Equal(response, gatewayResponse); diff != nil {
		t.Error(diff)
	}

}
func TestErrorRequestQuotation(t *testing.T) {
	initVariables()

	mockError := errors.New("MockError")
	lambdaHandler := awsHandler{
		reqSrv: mockRequestService{t: t, ctx: ctx, apiReq: apiRequest, reqQuotError: mockError},
		rspSrv: mockResponseService{t: t, ctx: ctx, requestId: gatewayRequestContext.RequestID, resp: &apiResponse},
	}

	gatewayResponse = events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}

	response, err := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.Equal(t, mockError, err, "Unexpected error")
	assert.NotNil(t, response, "Nil response")
	if diff := deep.Equal(response, gatewayResponse); diff != nil {
		t.Error(diff)
	}

}
func TestErrorResponseTimeout(t *testing.T) {
	initVariables()

	mockError := &service.TimeoutError{}
	lambdaHandler := awsHandler{
		reqSrv: mockRequestService{t: t, ctx: ctx, apiReq: apiRequest},
		rspSrv: mockResponseService{t: t, ctx: ctx, requestId: gatewayRequestContext.RequestID, respError: mockError},
	}

	gatewayResponse = events.APIGatewayProxyResponse{StatusCode: http.StatusRequestTimeout}

	response, err := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.Equal(t, mockError, err, "Unexpected error")
	assert.NotNil(t, response, "Nil response")
	if diff := deep.Equal(response, gatewayResponse); diff != nil {
		t.Error(diff)
	}

}
func TestErrorResponseError(t *testing.T) {
	initVariables()

	mockError := errors.New("Mock Error")
	lambdaHandler := awsHandler{
		reqSrv: mockRequestService{t: t, ctx: ctx, apiReq: apiRequest},
		rspSrv: mockResponseService{t: t, ctx: ctx, requestId: gatewayRequestContext.RequestID, respError: mockError},
	}

	gatewayResponse = events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}

	response, err := lambdaHandler.handleRequest(ctx, gatewayRequest)

	assert.Equal(t, mockError, err, "Unexpected error")
	assert.NotNil(t, response, "Nil response")
	if diff := deep.Equal(response, gatewayResponse); diff != nil {
		t.Error("Invalid Gateway Respose: ", diff)
	}

}

type mockRequestService struct {
	t            *testing.T
	ctx          context.Context
	apiReq       api.QuotationRequest
	reqQuotError error
}

func (rs mockRequestService) RequestQuotation(ctx context.Context, req api.QuotationRequest) error {
	assert.Equal(rs.t, rs.ctx, ctx, "Unexpected context")
	if diff := deep.Equal(rs.apiReq, req); diff != nil {
		rs.t.Error("Invalid QuotationRequest: ", diff)
	}
	return rs.reqQuotError
}

type mockResponseService struct {
	t         *testing.T
	ctx       context.Context
	requestId string
	resp      *api.QuotationResponse
	respError error
}

func (rs mockResponseService) WaitForResponse(ctx context.Context, requestId string) (*api.QuotationResponse, error) {
	assert.Equal(rs.t, rs.ctx, ctx, "Unexpected context")
	assert.Equal(rs.t, rs.requestId, requestId, "Unexpected requestId")

	return rs.resp, rs.respError
}
