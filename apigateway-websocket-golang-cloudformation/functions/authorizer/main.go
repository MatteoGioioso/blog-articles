package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go/v4"
)

func generatePolicy(principalId, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	// Optional output with custom properties of the String, Number or Boolean type.
	// This must be only primitive types
	authResponse.Context = map[string]interface{}{
		"stringKey":  "stringval",
		"numberKey":  123,
		"booleanKey": true,
	}
	return authResponse
}

func validateToken(token string) (*jwt.Token, error)  {
	decodedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// You can load your key from KMS here
		key := ""
		return key, nil
	})
	
	if err != nil {
		return &jwt.Token{}, err
	}
	
	return decodedToken, nil
}

func function(event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := event.QueryStringParameters["Auth"]

	_, err := validateToken(token)
	if err != nil {
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}
	
	return generatePolicy("user", "Allow", event.MethodArn), nil
}

func main() {
	lambda.Start(function)
}
