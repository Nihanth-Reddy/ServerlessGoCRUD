package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var CreatedUser CreateUserResponse

func generateRandomUser() Person {
	person = Person{
		FirstName: RandomString(6),
		LastName:  RandomString(6),
		Email:     RandomEmail(),
		Age:       int(RandomInt(0, 99)),
		Courses:   []string{"Test_New", "Test_New"},
	}
	return person
}

func CreateTestUser(t *testing.T) events.APIGatewayProxyResponse {
	person = generateRandomUser()
	json_user, _ := json.Marshal(person)
	req := events.APIGatewayProxyRequest{Body: string(json_user)}
	response, _ := CreateUser(req)
	if response.StatusCode != 201 {
		t.Errorf("Expected 201 got %d", response.StatusCode)
	}
	return response
}

func TestUpdateUser(t *testing.T) {
	person = generateRandomUser()
	json_user, _ := json.Marshal(person)
	req := events.APIGatewayProxyRequest{Body: string(json_user)}
	response, _ := CreateUser(req)
	if response.StatusCode != 201 {
		t.Errorf("While creating user Expected 201 got %d", response.StatusCode)
	}
	_ = json.Unmarshal([]byte(response.Body), &CreatedUser)

	person.Courses = []string{"Test_Updated", "Test_Updated"}
	json_user, _ = json.Marshal(person)
	req = events.APIGatewayProxyRequest{Body: string(json_user), QueryStringParameters: map[string]string{"user": CreatedUser.UserId}}
	defer DeleteUser(req)
	updateResponse, _ := UpdateUser(req)
	if updateResponse.StatusCode != 200 {
		t.Errorf("While creating user Expected 201 got %d", response.StatusCode)
	}
	var updatedPerson Person
	_ = json.Unmarshal([]byte(updateResponse.Body), &updatedPerson)
	if person.Courses[0] != updatedPerson.Courses[0] {
		t.Errorf("While pdating user expected courses are %v\ninstead got %v", person.Courses, updatedPerson.Courses)
	}
}

func TestCreateUser(t *testing.T) {
	createResponse := CreateTestUser(t)
	err := json.Unmarshal([]byte(createResponse.Body), &CreatedUser)
	if err != nil {
		t.Errorf("Unable to unmarshal created user response, ERROR: %v", err)
	}
	req := events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"user": CreatedUser.UserId}}
	delResponse, _ := DeleteUser(req)
	if delResponse.StatusCode != 204 {
		t.Logf("While deleting test user Expected 204, got %d", delResponse.StatusCode)
	}
}

func TestDeleteUser(t *testing.T) {
	TestCreateUser(t)
}

func TestFetchUser(t *testing.T) {
	createResponse := CreateTestUser(t)
	err := json.Unmarshal([]byte(createResponse.Body), &CreatedUser)
	if err != nil {
		t.Errorf("Unable to unmarshal created user response, ERROR: %v", err)
	}
	req := events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"user": CreatedUser.UserId}}
	defer DeleteUser(req)
	fetchResponse, err := FetchUser(req)
	if err != nil {
		t.Errorf("Unable to Fetch created user response, ERROR: %v", err)
	}
	if fetchResponse.StatusCode != 200 {
		t.Errorf("Fetch created user returned %d, Expected: 200", fetchResponse.StatusCode)
	}
	err = json.Unmarshal([]byte(fetchResponse.Body), &person)
	if err != nil {
		t.Errorf("Unable to unmarshal Fetch user response, ERROR: %v", err)
	}
	if CreatedUser.UserId != person.UserId {
		t.Errorf("Created User Id: %v is not matching with Fetched User Id: %v", CreatedUser.UserId, person.UserId)
	}
}

func FakeDBCall(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}

func TestHandler(t *testing.T) {
	d := Deps{
		Fu: FakeDBCall,
		Cu: FakeDBCall,
		Uu: FakeDBCall,
		Du: FakeDBCall,
	}
	_, err := d.handler(events.APIGatewayProxyRequest{HTTPMethod: "GET"})
	if err != nil {
		t.Errorf("For GET Call handler function is not supposed to return Error: %v", err)
	}
	_, err = d.handler(events.APIGatewayProxyRequest{HTTPMethod: "POST"})
	if err != nil {
		t.Errorf("For POST Call handler function is not supposed to return Error: %v", err)
	}
	_, err = d.handler(events.APIGatewayProxyRequest{HTTPMethod: "DELETE"})
	if err != nil {
		t.Errorf("For DELETE Call handler function is not supposed to return Error: %v", err)
	}
	_, err = d.handler(events.APIGatewayProxyRequest{HTTPMethod: "PUT"})
	if err != nil {
		t.Errorf("For PUT Call handler function is not supposed to return Error: %v", err)
	}
	res, err := d.handler(events.APIGatewayProxyRequest{HTTPMethod: "ANY"})
	if err != nil {
		t.Errorf("For UPDATE Call handler function is not supposed to return Error: %v", err)
	}
	if res.StatusCode != 405 {
		t.Errorf("For ANY Call handler function is supposed to return 405: instead it returned %d", res.StatusCode)
	}

}
