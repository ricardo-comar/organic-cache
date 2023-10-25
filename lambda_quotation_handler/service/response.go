package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ricardo-comar/organic-cache/lib_common/api"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/quotation_handler/gateway"
)

type TimeoutError struct {
}

func (*TimeoutError) Error() string {
	return "Timeout waiting for response"
}

var timeout = 10 * time.Second

type MessageChannel struct {
	Resp  api.QuotationResponse
	Error error
}

func WaitForResponse(ctx context.Context, sqscli *sqs.Client, requestId string) (*api.QuotationResponse, error) {

	response := make(chan MessageChannel)

	go waitForMessage(response, ctx, sqscli, requestId)

	select {

	case resp := <-response:
		log.Println("Processamento finalizado")
		return &resp.Resp, resp.Error

	case <-time.After(timeout):
		log.Println("Timeout! A tarefa demorou muito para ser concluída.")
		return nil, &TimeoutError{}
	}

}

func waitForMessage(response chan MessageChannel, ctx context.Context, sqscli *sqs.Client, requestId string) {

	for {

		result, err := gateway.ReceiveMessage(ctx, sqscli)

		if err != nil {
			log.Println("Erro ao receber mensagem:", err)
			response <- MessageChannel{Error: err}
			return
		}

		for _, msg := range result.Messages {

			if reqId, found := msg.MessageAttributes["RequestId"]; found && requestId == *reqId.StringValue {
				log.Println("Mensagem relevante:", *msg.Body)

				quotation := message.QuotationMessage{}
				json.Unmarshal([]byte(*msg.Body), &quotation)
				if err != nil {
					log.Println("Erro ao transformar mensagem em struct:", err)
					response <- MessageChannel{Error: err}
					return
				}

				log.Println("Removendo a mensagem:", msg.ReceiptHandle)
				gateway.DeleteMessage(ctx, sqscli, msg.ReceiptHandle)

				if err != nil {
					log.Println("Erro ao remover a mensagem:", err)
					response <- MessageChannel{Error: err}
					return
				}
				resp := api.QuotationResponse{
					RequestId: quotation.RequestId,
					UserId:    quotation.UserId,
					Products:  quotation.Products,
				}
				log.Printf("Resposta: %+v", resp)

				response <- MessageChannel{Resp: resp}
				return

			} else {

				err = gateway.ChangeMessageVisibility(ctx, sqscli, msg.ReceiptHandle)

				if err != nil {
					log.Println("Erro ao devolver a mensagem à fila:", err)
					response <- MessageChannel{Error: err}
					return
				}
			}

		}

	}
}
