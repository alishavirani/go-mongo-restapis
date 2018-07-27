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
		GetEmployees,
	},
	Route{
		"AddEmployee",
		"POST",
		"/employee",
		AddEmployee,
	},
	Route{
		"GetEmployeeRecord",
		"GET",
		"/employee/{id}",
		GetEmployeeRecord,
	},
	Route{
		"UpdateEmployee",
		"PUT",
		"/employee",
		UpdateEmployee,
	}, Route{
		"DeleteEmployee",
		"DELETE",
		"/employee/{id}",
		DeleteEmployee,
	},
}
