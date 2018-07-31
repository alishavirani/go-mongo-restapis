package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")

		type Exception struct {
			Message string
		}

		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			fmt.Println("Bearer token?", bearerToken)

			if len(bearerToken) == 1 {
				token, error := jwt.Parse(bearerToken[0], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					log.Println("TOKEN WAS VALID")
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

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

	fmt.Println("emp in handler???", employee)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if (Employee{}) == employee {
		fmt.Println("No employee record found")
		w.Write([]byte("No employee record found"))
		return
	}

	data, _ := json.Marshal(employee)
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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if updated {
		data := []byte("Employee Updated Successfully")
		w.Write(data)
	} else {
		data := []byte("No such employee found")
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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if deleted {
		data := []byte("Deleted successfully")
		w.Write(data)
	} else {
		data := []byte("No such employee found")
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
		fmt.Println("Error in umnarshal???", err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("Error in json encoder???", err)
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if result := AddEmployeeToDb(employee); result {
		fmt.Println("Added record to DB")
		w.Write([]byte("Added record to DB"))
	} else {
		fmt.Println("Failed to insert to DB")
		fmt.Fprintf(w, "failed to insert to DB")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// func RegisterEmployee(w http.ResponseWriter, r *http.Request) {
// 	var user UserAccess
// 	fmt.Println("REQ???", r)
// 	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

// 	fmt.Println("Body???", body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err := r.Body.Close(); err != nil {
// 		panic(err)
// 	}
// 	if err := json.Unmarshal(body, &user); err != nil {
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(422) // unprocessable entity
// 		if err := json.NewEncoder(w).Encode(err); err != nil {
// 			panic(err)
// 		}
// 	}

// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	fmt.Println("User???", user)

// 	// _ = json.NewDecoder(r.Body).Decode(&user)

// 	// fmt.Println("User#######", user)

// 	if result := RegisterEmployeeToDb(user); result {
// 		fmt.Println("User result???", result)
// 		fmt.Println("Registered User In Db")
// 		w.Write([]byte("Registered User In Db"))
// 	} else {
// 		fmt.Println("Failed to insert to DB")
// 		fmt.Fprintf(w, "failed to insert to DB")
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// }

func RegisterEmployee(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user UserAccess
		fmt.Println("REQ???", r)
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

		fmt.Println("Body???", r.Body)
		fmt.Println("Body byte???", body)
		if err != nil {
			panic(err)
		}
		// if err := r.Body.Close(); err != nil {
		// 	panic(err)
		// }
		if err := json.Unmarshal(body, &user); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		fmt.Println("User???", user)

		if result := RegisterEmployeeToDb(user); result {
			fmt.Println("Registered User In Db")
			// w.Write([]byte("Registered User In Db"))
			next(w, r)
		} else {
			fmt.Println("Failed to insert to DB")
			fmt.Fprintf(w, "failed to insert to DB")
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
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
		fmt.Println("No such user found")
		fmt.Fprintf(w, "No such user found")
	} else {
		fmt.Println("Logged in successfully")
		w.Write([]byte("Logged in successfully"))
	}
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	var user UserAccess

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	fmt.Println("Body??? token", r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Body??? byte", body)
	// _ = json.NewDecoder(r.Body).Decode(&user)
	// fmt.Print("User###########", user)

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

	fmt.Println("User???#####", user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	type JwtToken struct {
		Token string
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}
