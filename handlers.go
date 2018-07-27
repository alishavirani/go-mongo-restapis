package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := FetchAllEmployeesFromDb()

	if err != nil {
		panic(err)
	}

	data, _ := json.Marshal(employees)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetEmployeeRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("ID???", id)
	employee, err := GetEmployeeById(id)

	if err != nil {
		panic(err)
	}

	data, _ := json.Marshal(employee)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &employee); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	updated, err := UpdateEmployeeById(employee)

	if err != nil {
		panic(err)
	}
	if updated {
		data := []byte("Employee Updated Successfully")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}

}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("ID???", id)
	deleted, err := DeleteEmployeeById(id)

	if err != nil {
		panic(err)
	}

	if deleted {
		data := []byte("Deleted successfully")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}

}

func AddEmployee(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &employee); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if result := AddEmployeeToDb(employee); result {
		fmt.Println("Added record to DB")
		w.Write([]byte("Failed to insert to DB"))
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("Failed to insert to DB")
		w.Write([]byte("Failed to insert to DB"))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func RegisterEmployee(w http.ResponseWriter, r *http.Request) {
	var user UserAccess
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if result := RegisterEmployeeToDb(user); result {
		fmt.Println("Registered User In Db")
		w.Write([]byte("Registered User In Db"))
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("Failed to insert to DB")
		w.Write([]byte("Failed to insert to DB"))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user UserAccess
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	result, err := LoginDb(user)

	if err != nil {
		panic(err)
	}
	fmt.Print("Result????", result)
	if result.Email == "" {
		fmt.Println("Registered User In Db")
		w.Write([]byte("Registered User In Db"))
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("Logged in successfully")
		w.Write([]byte("Logged in successfully"))
		w.WriteHeader(http.StatusOK)
	}
}
