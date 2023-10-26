package service

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/ricardo-comar/organic-cache/lib_common/api"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
	"github.com/ricardo-comar/organic-cache/quotation_handler/gateway"
)

type TimeoutError struct {
}

func (*TimeoutError) Error() string {
	return "Timeout waiting for response"
}

type MessageChannel struct {
	Resp  *api.QuotationResponse
	Error error
}

type responseService struct {
	msgWait time.Duration
	sqsg    gateway.SQSGateway
}

func NewResponseService(msgWait time.Duration, sqsg gateway.SQSGateway) ResponseService {
	rs := &responseService{
		msgWait: msgWait,
		sqsg:    sqsg,
	}
	return ResponseService(rs)
}

type ResponseService interface {
	WaitForResponse(ctx context.Context, requestId string) (*api.QuotationResponse, error)
}

func (rs responseService) WaitForResponse(ctx context.Context, requestId string) (*api.QuotationResponse, error) {

	response := make(chan MessageChannel)

	go waitForMessage(response, ctx, rs.sqsg, requestId)

	select {

	case resp := <-response:
		log.Println("Processamento finalizado")
		return resp.Resp, resp.Error

	case <-time.After(rs.msgWait):
		log.Println("Timeout! A tarefa demorou muito para ser concluída.")
		return nil, &TimeoutError{}
	}

}

func waitForMessage(response chan MessageChannel, ctx context.Context, sqsg gateway.SQSGateway, requestId string) {

	for {

		result, err := sqsg.ReceiveMessage(ctx)

		if err != nil {
			log.Println("Erro ao receber mensagem:", err)
			response <- MessageChannel{Error: err}
			return
		}

		for _, msg := range result.Messages {

			if reqId, found := msg.MessageAttributes["RequestId"]; found && requestId == *reqId.StringValue {
				log.Println("Mensagem relevante:", *msg.Body)

				quotation := message.QuotationMessage{}
				decoder := json.NewDecoder(strings.NewReader(*msg.Body))
				decoder.DisallowUnknownFields()
				err = decoder.Decode(&quotation)
				if err != nil {
					log.Println("Erro ao transformar mensagem em struct:", err)
					response <- MessageChannel{Error: err}
					return
				}

				log.Println("Removendo a mensagem:", msg.ReceiptHandle)
				err = sqsg.DeleteMessage(ctx, msg.ReceiptHandle)
				if err != nil {
					log.Println("Erro ao remover a mensagem:", err)
					response <- MessageChannel{Error: err}
					return
				}

				resp := &api.QuotationResponse{
					RequestId: quotation.RequestId,
					UserId:    quotation.UserId,
					Products:  quotation.Products,
				}
				log.Printf("Resposta: %+v", resp)

				response <- MessageChannel{Resp: resp}
				return

			} else {

				err = sqsg.ChangeMessageVisibility(ctx, msg.ReceiptHandle)

				if err != nil {
					log.Println("Erro ao devolver a mensagem à fila:", err)
					response <- MessageChannel{Error: err}
					return
				}
			}

		}

	}
}
