package main

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const SERVER = "mongodb://localhost/go-rest-api"

const DBNAME = "go-rest-api"

const COLLECTION = "employees"

func AddEmployeeToDb(employee Employee) bool {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer session.Close()

	session.DB(DBNAME).C(COLLECTION).Insert(employee)
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

	c := session.DB(DBNAME).C(COLLECTION)

	results := Employees{}

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
		return nil, err
	}
	fmt.Println("printing db docs!!!!", results)
	return results, nil

}

func GetEmployeeById(id string) (Employees, error) {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	result := Employees{}

	filter := bson.M{"id": id}

	fmt.Println("filter????", filter)

	if err := c.Find(&filter).Limit(5).All(&result); err != nil {
		fmt.Println("Failed to write result:", err)
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

	err = session.DB(DBNAME).C(COLLECTION).Update(&filter, employee)

	if err != nil {
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

	filter := bson.M{"id": "101"}

	fmt.Println("Filter????", filter)

	if err = session.DB(DBNAME).C(COLLECTION).Remove(&filter); err != nil {
		log.Fatal(err)
		return false, err
	}

	fmt.Println("Deleted Product ID - ", id)
	return true, nil

}
