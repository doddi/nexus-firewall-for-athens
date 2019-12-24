package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"nexus-firewall-for-athens/athens"
	"nexus-firewall-for-athens/validate"
)

type LambdaServer struct {
	Validator validate.Validator
}

func (server LambdaServer) Handle() {
	server.handleLambda()
}

func (server LambdaServer) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var buffer bytes.Buffer

	fmt.Println("Request received")

	hookMessage := server.decodeMessage(request.Body)
	_ = server.Validator.Validate(hookMessage)

	return events.APIGatewayProxyResponse{
		Body:            base64.StdEncoding.EncodeToString(buffer.Bytes()),
		StatusCode:      200,
		IsBase64Encoded: true,
	}, nil
}

func (server LambdaServer) decodeMessage(request string) athens.Request {
	hookMessage := athens.Request{}
	json.Unmarshal([]byte(request), &hookMessage)
	return hookMessage
}

func (server LambdaServer) handleLambda() {
	lambda.Start(server.handler)
}
