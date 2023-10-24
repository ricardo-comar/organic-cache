package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ricardo-comar/organic-cache/lib_common/api"
	"github.com/ricardo-comar/organic-cache/lib_common/message"
)

type TimeoutError struct {
}

func (*TimeoutError) Error() string {
	return "Timeout waiting for response"
}

var queueURL = os.Getenv("QUOTATIONS_QUEUE")
var timeout = 10 * time.Second

func WaitForResponse(ctx context.Context, sqscli *sqs.Client, requestId string) (*api.QuotationResponse, error) {

	taskCompleted := make(chan bool)
	response := make(chan *api.QuotationResponse)
	err := make(chan error)

	go waitForMessage(taskCompleted, response, err, ctx, sqscli, requestId)

	select {

	case <-taskCompleted:
		fmt.Println("Processamento finalizado")
		return <-response, <-err

	case <-time.After(timeout):
		fmt.Println("Timeout! A tarefa demorou muito para ser concluída.")
		return nil, &TimeoutError{}
	}

}

func waitForMessage(done chan bool, response chan *api.QuotationResponse, errChan chan error, ctx context.Context, sqscli *sqs.Client, requestId string) {

	for {
		result, err := sqscli.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:              &queueURL,
			MaxNumberOfMessages:   1,
			WaitTimeSeconds:       3,
			AttributeNames:        []types.QueueAttributeName{"All"},
			MessageAttributeNames: []string{"All"},
		})

		if err != nil {
			fmt.Println("Erro ao receber mensagem:", err)
			errChan <- err
			done <- true
		}

		for _, msg := range result.Messages {
			request_id := msg.MessageAttributes["RequestId"].StringValue

			if requestId, found := msg.MessageAttributes["RequestId"]; found && *request_id == *requestId.StringValue {

				fmt.Println("Mensagem relevante:", *msg.Body)
				quotation := message.QuotationMessage{}
				json.Unmarshal([]byte(*msg.Body), &quotation)
				if err != nil {
					fmt.Println("Erro ao transformar mensagem em struct:", err)
					errChan <- err
					done <- true
				}

				_, err := sqscli.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &queueURL,
					ReceiptHandle: msg.ReceiptHandle,
				})
				if err != nil {
					fmt.Println("Erro ao excluir a mensagem:", err)
					errChan <- err
					done <- true
				}

				response <- &api.QuotationResponse{
					RequestId: quotation.RequestId,
					UserId:    quotation.UserId,
					Products:  quotation.Products,
				}
				done <- true
			} else {

				_, err := sqscli.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
					QueueUrl:          &queueURL,
					ReceiptHandle:     msg.ReceiptHandle,
					VisibilityTimeout: 0, // Defina para 0 para devolver imediatamente à fila
				})
				if err != nil {
					fmt.Println("Erro ao devolver a mensagem à fila:", err)
				}
			}

		}

	}
}
