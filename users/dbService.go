package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

type Person struct {
	UserId    string   `json:"user_id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Age       int      `json:"age"`
	Courses   []string `json:"courses"`
}

type CreateUserResponse struct {
	UserId string `json:"user_id"`
}

const TableName = "GolangUsers"

var person Person

func CreateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{}
	dbClient := getDBClient()
	err := json.Unmarshal([]byte(req.Body), &person)
	if err != nil {
		return errorResponse(response, 400, "Recieved an error while unmarshalling the request body: %v\n", err)
	}

	person.UserId = uuid.NewString()
	log.Printf("Person Object for DB marshalling %#v\n", person)
	attributeValue, err := dynamodbattribute.MarshalMap(person)
	if err != nil {
		return errorResponse(response, 500, "Recieved an error while Marshalling the person struct: %v\n", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(TableName),
	}
	_, err = dbClient.PutItem(input)

	if err != nil {
		return errorResponse(response, 500, "Got error calling PutItem: %s\n", err)
	}
	response.StatusCode = 201
	response.Headers = map[string]string{"Content-Type": "application/json"}

	res := CreateUserResponse{UserId: person.UserId}
	res_str, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return errorResponse(response, 500, "Recieved an error while marshalling the response body: %v\n", err)
	}
	response.Body = string(res_str)
	log.Printf("Response of the create user call is %#v\n", response)
	return response, nil
}

func FetchUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{}
	response.Headers = map[string]string{"Content-Type": "application/json"}
	dbClient := getDBClient()
	UserId := req.QueryStringParameters["user"]
	log.Printf("Searching for user with id: %s", UserId)

	item := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(UserId),
			},
		},
	}
	res, err := dbClient.GetItem(item)
	if err != nil {
		return errorResponse(response, 500, "Got error calling GetItem: %s", err)
	}

	if res.Item == nil {
		return errorResponse(response, 404, "Couldn't get the users with id: "+UserId+"\n", nil)
	}

	err = dynamodbattribute.UnmarshalMap(res.Item, &person)
	if err != nil {
		return errorResponse(response, 500, "Got error while Unmarshalling response: %s", err)
	}

	res_str, err := json.MarshalIndent(person, "", " ")
	if err != nil {
		return errorResponse(response, 500, "Got error while marshalling response: %s", err)
	}

	response.StatusCode = 200
	response.Body = string(res_str)

	return response, nil
}

func UpdateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{}
	response.Headers = map[string]string{"Content-Type": "application/json"}
	dbClient := getDBClient()
	UserId := req.QueryStringParameters["user"]
	log.Printf("Updating user with id: %s", UserId)

	err := json.Unmarshal([]byte(req.Body), &person)
	if err != nil {
		return errorResponse(response, 400, "Recieved an error while unmarshalling the request body: %v\n", err)
	}

	person.UserId = UserId

	log.Printf("Person Object for DB marshalling %#v\n", person)
	attributeValue, err := dynamodbattribute.MarshalMap(person)
	if err != nil {
		return errorResponse(response, 500, "Recieved an error while Marshalling the person struct: %v\n", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(TableName),
	}
	_, err = dbClient.PutItem(input)

	if err != nil {
		return errorResponse(response, 500, "Got error calling PutItem: %s", err)
	}

	res_str, err := json.MarshalIndent(person, "", " ")
	if err != nil {
		return errorResponse(response, 500, "Got error while marshalling response: %s", err)
	}

	response.StatusCode = 200
	response.Body = string(res_str)
	log.Printf("Response of the update user call is %#v\n", response)
	return response, nil
}

func DeleteUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{}
	response.Headers = map[string]string{"Content-Type": "application/json"}
	dbClient := getDBClient()

	UserId := req.QueryStringParameters["user"]

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(UserId),
			},
		},
	}

	_, err := dbClient.DeleteItem(input)
	if err != nil {
		return errorResponse(response, 500, "Got error while deleting the item from DB: %v", err)
	}
	response.StatusCode = 204
	return response, nil
}

func getDBClient() dynamodbiface.DynamoDBAPI {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ServiceClient := dynamodb.New(sess)
	return ServiceClient
}

func errorResponse(res events.APIGatewayProxyResponse, code int, message string, err error) (events.APIGatewayProxyResponse, error) {
	res.StatusCode = code
	log.Printf(message, err)
	return res, nil
}
