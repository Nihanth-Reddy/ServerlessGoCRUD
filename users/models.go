package main

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
