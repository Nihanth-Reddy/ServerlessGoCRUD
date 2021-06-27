package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Incoming request event is: %#v\n", req)
	switch req.HTTPMethod {
	case "GET":
		return FetchUser(req)
	case "POST":
		return CreateUser(req)
	case "PUT":
		return UpdateUser(req)
	case "DELETE":
		return DeleteUser(req)
	default:
		log.Printf("Requested for unhandled method %#v", req.HTTPMethod)
		response := events.APIGatewayProxyResponse{StatusCode: 405}
		return response, nil
	}
}
