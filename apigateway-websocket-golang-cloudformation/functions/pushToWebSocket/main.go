package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"os"
)

type RequestBody struct {
	RepositoryUrl string `json:"repositoryUrl"`
}

var (
	requestBody        = RequestBody{}
	initializedSession = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	apiGatewayId = os.Getenv("API_GATEWAY_ID")
	environment  = os.Getenv("GO_ENV")
	region       = os.Getenv("AWS_REGION")
	api          = NewApiGatewayManagementApi()
)

func GetApiGatewayEndpoint(apiGatewayId string) string {
	return fmt.Sprintf("%v.execute-api.%v.amazonaws.com/%v", apiGatewayId, region, environment)
}

func NewApiGatewayManagementApi() *apigatewaymanagementapi.ApiGatewayManagementApi {
	return apigatewaymanagementapi.New(initializedSession,
		aws.NewConfig().WithEndpoint(GetApiGatewayEndpoint(apiGatewayId)))
}

func function(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := request.RequestContext.ConnectionID
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Something went wrong",
		}, nil
	}

	input := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionId),
		Data:         []byte("Hello there!"),
	}
	if _, err := api.PostToConnection(input); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Something went wrong",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(function)
}
