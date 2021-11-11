package main

import (
	"encoding/json"
	"fmt"
	"liturgia/index/funcao"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/net/context"
)

type Requisicao struct {
	Data time.Time `json:"data"`
}

func handlerRequest(ctx context.Context, requisicao Requisicao) (events.APIGatewayProxyResponse, error) {
	fmt.Println(requisicao)
	obj := funcao.BuscarLiturgia(requisicao.Data)
	body, err := json.Marshal(obj)
	// ctx.
	// jsonObj := json.end
	// return jsonObj, nil
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handlerRequest)
	// funcao.BuscarLiturgia(time.Now())
}
