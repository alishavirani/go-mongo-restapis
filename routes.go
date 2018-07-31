package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"GetEmployees",
		"GET",
		"/employeeList",
		AuthenticationMiddleware(GetEmployees),
	},
	Route{
		"AddEmployee",
		"POST",
		"/employee",
		AuthenticationMiddleware(AddEmployee),
	},
	Route{
		"GetEmployeeRecord",
		"GET",
		"/employee/{id}",
		AuthenticationMiddleware(GetEmployeeRecord),
	},
	Route{
		"UpdateEmployee",
		"PUT",
		"/employee",
		AuthenticationMiddleware(UpdateEmployee),
	},
	Route{
		"DeleteEmployee",
		"DELETE",
		"/employee/{id}",
		AuthenticationMiddleware(DeleteEmployee),
	},
	Route{
		"RegisterEmployee",
		"POST",
		"/register",
		RegisterEmployee(GetToken),
	},
	Route{
		"LoginUser",
		"POST",
		"/login",
		AuthenticationMiddleware(LoginUser),
	},
	// Route{
	// 	"Authentication",
	// 	"POST",
	// 	"/get-token",
	// 	GetToken,
	// },
}
