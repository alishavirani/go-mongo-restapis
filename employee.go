package main

type Employee struct {
	Name       string `json: "name"`
	Email      string `json: "email"`
	ID         string `json: "id`
	Phone      uint64 `json: "phone"`
	Address    string `json: "address"`
	Department string `json: "department"`
}

type Employees []Employee
