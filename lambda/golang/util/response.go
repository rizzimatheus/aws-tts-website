package util

import "github.com/aws/aws-lambda-go/events"

// NewAPIGatewayProxyResponse creates a new API Gateway Proxy Response.
func NewAPIGatewayProxyResponse(status int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Methods":     "GET, POST, OPTIONS",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: body,
	}
}
