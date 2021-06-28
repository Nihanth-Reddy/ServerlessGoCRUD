package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type helperFunc func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Deps struct {
	Fu helperFunc
	Cu helperFunc
	Uu helperFunc
	Du helperFunc
}

func main() {
	d := Deps{}
	lambda.Start(d.handler)
}

func (d *Deps) handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Incoming request event is: %#v\n", req)
	switch req.HTTPMethod {
	case "GET":
		if d.Fu == nil {
			return FetchUser(req)
		}
		return d.Fu(req)
	case "POST":
		if d.Cu == nil {
			return CreateUser(req)
		}
		return d.Cu(req)
	case "PUT":
		if d.Uu == nil {
			return UpdateUser(req)
		}
		return d.Uu(req)
	case "DELETE":
		if d.Du == nil {
			return DeleteUser(req)
		}
		return d.Du(req)
	default:
		log.Printf("Requested for unhandled method %#v", req.HTTPMethod)
		response := events.APIGatewayProxyResponse{StatusCode: 405}
		return response, nil
	}
}
