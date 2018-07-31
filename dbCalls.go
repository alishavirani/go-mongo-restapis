package main

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const SERVER = "mongodb://localhost/go-rest-api"

const DBNAME = "go-rest-api"

const EMPCOLLECTION = "employees"

const ACCESSCOLLECTION = "access"

func RegisterEmployeeToDb(user UserAccess) bool {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer session.Close()

	session.DB(DBNAME).C(ACCESSCOLLECTION).Insert(user)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Registered new emp- ", user)
	return true
}

func LoginDb(user UserAccess) (UserInDb UserAccess, err error) {
	var result UserAccess
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
		return result, err
	}
	defer session.Close()

	c := session.DB(DBNAME).C(ACCESSCOLLECTION)

	filter := bson.M{"email": user.Email}

	if err := c.Find(&filter).One(&result); err != nil {
		return result, nil
	}
	return result, nil
}

func AddEmployeeToDb(employee Employee) bool {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer session.Close()

	session.DB(DBNAME).C(EMPCOLLECTION).Insert(employee)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Added New Employee ID- ", employee)
	return true
}

func FetchAllEmployeesFromDb() (Employees, error) {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer session.Close()

	c := session.DB(DBNAME).C(EMPCOLLECTION)

	results := Employees{}

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
		return nil, err
	}
	fmt.Println("printing db docs!!!!", results)
	return results, nil

}

func GetEmployeeById(id string) (Employee, error) {
	session, err := mgo.Dial(SERVER)

	var result Employee

	if err != nil {
		log.Fatal(err)
		return result, err
	}
	defer session.Close()

	c := session.DB(DBNAME).C(EMPCOLLECTION)

	filter := bson.M{"id": id}

	if err := c.Find(&filter).One(&result); err != nil {
		fmt.Println("printing error in db", err)

		if err.Error() == "not found" {
			return result, nil
		}
		log.Fatal(err)
		return result, err
	}
	return result, nil
}

func UpdateEmployeeById(employee Employee) (bool, error) {

	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer session.Close()

	filter := bson.M{"id": employee.ID}

	if err = session.DB(DBNAME).C(EMPCOLLECTION).Update(&filter, employee); err != nil {
		if err.Error() == "not found" {
			return false, nil
		}
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func DeleteEmployeeById(id string) (bool, error) {

	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer session.Close()

	filter := bson.M{"id": id}

	if err = session.DB(DBNAME).C(EMPCOLLECTION).Remove(&filter); err != nil {

		if err.Error() == "not found" {
			return false, nil
		}
		log.Fatal(err)
		return false, err
	}

	fmt.Println("Deleted Product ID - ", id)
	return true, nil

}
