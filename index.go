package main

import (
	"encoding/json"
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
	// fmt.Println(requisicao)
	obj := funcao.BuscarLiturgia(requisicao.Data)
	body, err := json.Marshal(obj)
	// ctx.
	// jsonObj := json.end
	// return jsonObj, nil
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200, Headers: map[string]string{"Content-Type": "application/json"}}, nil
}

func main() {
	lambda.Start(handlerRequest)
	// dataSel := time.Now()
	// dataSel, err := time.Parse("2006-01-02", "2021-12-25")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// funcao.BuscarLiturgia(dataSel)
}
