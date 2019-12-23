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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var buffer bytes.Buffer

	fmt.Println("Request received")

	hookMessage := decodeLambdaMessage(request.Body)
	_ = validate.Validate(hookMessage)

	return events.APIGatewayProxyResponse{
		Body:            base64.StdEncoding.EncodeToString(buffer.Bytes()),
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/pdf"},
		IsBase64Encoded: true,
	}, nil
}

func decodeLambdaMessage(request string) athens.Request {
	hookMessage := athens.Request{}
	json.Unmarshal([]byte(request), &hookMessage)
	return hookMessage
}

func HandleLambda() {
	lambda.Start(Handler)
}
