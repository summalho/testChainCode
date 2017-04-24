package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

type SimpleChaincode struct {
}

func main() {

	fmt.Println("Inside main method")
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)

	}
}

type User struct {
	FfId       string
	Title      string
	Gender     string
	FirstName  string
	LastName   string
	Dob        string
	Email      string
	Country    string
	Address    string
	City       string
	Zip        string
	CreatedBy  string
	TotalPoint string
}

func (t *SimpleChaincode) registerUser(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	var newUserByteArray []byte

	newUser := User{}

	newUser.FfId = args[0]
	newUser.Title = args[1]
	newUser.Gender = args[2]
	newUser.FirstName = args[3]
	newUser.LastName = args[4]
	newUser.Dob = args[5]
	newUser.Email = args[6]
	newUser.Address = args[7]
	newUser.City = args[8]
	newUser.Zip = args[9]
	newUser.Country = args[10]
	newUser.CreatedBy = args[11]
	newUser.TotalPoint = args[12]

	// Marshal newUser to convert it into bytes and store in Block chain.
	newUserByteArray, err = json.Marshal(newUser)

	err = stub.PutState(newUser.FfId, newUserByteArray)

	if err != nil {
		fmt.Println("Could not save userDetails to ledger", err)
		return nil, err

	}

	return newUserByteArray, nil

}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	var bytes []byte

	if len(args) == 13 {

		bytes, err = t.registerUser(stub, function, args)

		if err != nil {
			return nil, err
		}

		return bytes, nil

	}

	return nil, errors.New("Add 13 arguments")

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	var newPoints int
	var userDataBytes []byte
	var points int

	if len(args) != 3 {
		return nil, errors.New("must have two arguments")
	}

	var userId = args[0]

	existingUser := User{}
	userDataBytes, err = stub.GetState(userId)

	if err != nil {
		return nil, errors.New("Error retreiving User with given Id")
	}

	err = json.Unmarshal(userDataBytes, &existingUser)

	if err != nil {
		return nil, errors.New("Problem in unmarshalling data")
	}

	pointsrecieved := existingUser.TotalPoint

	points, err = strconv.Atoi(pointsrecieved)

	if function == "addDeletePoints" {

		fmt.Println("Inside addDeletePoints")

		newPoints = t.addDeletePoints(stub, args, points)

	}

	// After addition or subtraction of points,  Store points back to the ledger.

	var finalPointsStr = strconv.Itoa(newPoints)
	existingUser.TotalPoint = finalPointsStr
	userDataBytes, err = json.Marshal(&existingUser)

	err = stub.PutState(userId, userDataBytes)

	return userDataBytes, nil

}

func (t *SimpleChaincode) addDeletePoints(stub shim.ChaincodeStubInterface, args []string, points int) int {

	var numberOfPointsToAddOrSub = args[1]
	numberOfPointsInt, err := strconv.Atoi(numberOfPointsToAddOrSub)

	if err != nil {

		return 0
	}

	addOrDelete := args[2]

	fmt.Println(addOrDelete, " = fmt.Println(addOrDelete)")

	if addOrDelete == "add" {
		points = points + numberOfPointsInt
	}

	if addOrDelete == "delete" {
		points = points - numberOfPointsInt

	}

	return points

}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var currentPoints []byte
	var userdetails []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("must have only 1 arguments")
	}

	if function == "getPoints" {
		currentPoints, err = t.getPoints(stub, args)
	}
	if function == "getUser" {
		userdetails, err = t.getUser(stub, args)
	}
	if err != nil {
		return nil, errors.New("returning from query with error. Current points or User Details are not correct.")
	}

	if function == "getPoints" {
		return currentPoints, nil
	}
	if function == "getUser" {
		return userdetails, nil
	}

	return nil, errors.New("getUser or getPoints did not worked properly")

}

func (t *SimpleChaincode) getUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userId = args[0]

	bytes, err := stub.GetState(userId)

	if err != nil {

		jsonResp := "{\"Error\":\"Failed to get state for " + userId + "\"}"
		return nil, errors.New(jsonResp)

	}

	userRetrieved := User{}
	err = json.Unmarshal(bytes, &userRetrieved)

	jsonResp := []byte(" {\"TotalPoints\":\"" + userRetrieved.TotalPoint + "\"}" + " {\"FirstName\":\"" + userRetrieved.FirstName + "\"}" + " {\"LastName\":\"" + userRetrieved.LastName + "\"}" + " {\"Id\":\"" + userRetrieved.FfId + "\"}" + " {\"DOB\":\"" + userRetrieved.Dob + "\"}")

	return jsonResp, nil

}

func (t *SimpleChaincode) getPoints(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var userId = args[0]

	bytes, err := stub.GetState(userId)

	if err != nil {

		jsonResp := "{\"Error\":\"Failed to get state for " + userId + "\"}"
		return nil, errors.New(jsonResp)

	}

	userRetrieved := User{}
	err = json.Unmarshal(bytes, &userRetrieved)

	jsonResp := []byte("{\"TotalPoints\":\"" + userRetrieved.TotalPoint + "\"}")

	return jsonResp, nil

}
