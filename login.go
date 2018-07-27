package main

type UserAccess struct {
	Email    string `json: "email"`
	Password string `json: password`
}

type Login []UserAccess
