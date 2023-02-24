package gateway

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"os"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/lambda"
// 	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
// )

// func RequestPriceCalc(ctx context.Context, lamcli lambda.Client, userId string) error {

// 	payload, err := json.Marshal(map[string]interface{}{
// 		"id": userId,
// 	})

// 	input := &lambda.InvokeInput{
// 		FunctionName:   aws.String(os.Getenv("PRICE_CALC_LAMBDA")),
// 		InvocationType: types.InvocationTypeEvent,
// 		Payload:        payload,
// 		Qualifier:      aws.String("1"),
// 	}
// 	log.Printf("Lambda %s input: %+v", os.Getenv("PRICE_CALC_LAMBDA"), input)

// 	result, err := lamcli.Invoke(ctx, input)

// 	log.Printf("Lambda response: %+v", result)
// 	log.Printf("Lambda error: %+v", err)

// 	return err
// }
