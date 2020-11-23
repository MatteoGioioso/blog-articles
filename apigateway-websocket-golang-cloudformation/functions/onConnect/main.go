package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// This function will be triggered once the Websocket connect
func function(event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := event.RequestContext.ConnectionID
	fmt.Println("Connection id: ", id)
	return events.APIGatewayProxyResponse{
		StatusCode:        200,
	}, nil
}

func main() {
	lambda.Start(function)
}
