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
	ffId       int         `json:"ffid"`
	title      string      `json:"title"`
	gender     string      `json:"gender"`
	firstName  string      `json:"firstname"`
	lastName   string      `json:"lastname"`
	dob        string      `json:"dob"`
	email      string      `json:"email"`
	userAdd    UserAddress `json:"userAdd"`
	createdBy  string      `json:"createdBy"`
	totalPoint int         `json:"totalPoint"`
}

type UserAddress struct {
	country     string `json:"country"`
	addressLine string `json:"addressLine"`
	city        string `json:"city"`
	zip         string `json:"zip"`
}

func (t *SimpleChaincode) registerNewUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var newUser User
	var newUserByteArray []byte
	var err error
	var id, tPoint int

	if len(args) != 13 {
		return nil, errors.New("must have thirteen arguments")
	}

	var idString string
	var tpoint string

	idString = args[0]
	tpoint = args[12]

	id, err = strconv.Atoi(idString)
	tPoint, err = strconv.Atoi(tpoint)

	newUser.ffId = id
	newUser.title = args[1]
	newUser.gender = args[2]
	newUser.firstName = args[3]
	newUser.lastName = args[4]
	newUser.dob = args[5]
	newUser.email = args[6]
	newUser.userAdd.addressLine = args[7]
	newUser.userAdd.city = args[8]
	newUser.userAdd.zip = args[9]
	newUser.userAdd.country = args[10]
	newUser.createdBy = args[11]
	newUser.totalPoint = tPoint

	var newUserDataAddress = &newUser
	newUserByteArray, err = json.Marshal(newUserDataAddress)

	err = stub.PutState(idString, newUserByteArray)

	if err != nil {
		return nil, errors.New("data cannot be pushed successfully. returned from registerNewUser")
	}

	return nil, nil

}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("Inside Init")

	if function == "registerNewUser" {
		return t.registerNewUser(stub, args)
	}

	return nil, errors.New("new user cannot be added successfully, returned from INIT")

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var userDataBytes []byte

	if len(args) != 2 {
		return nil, errors.New("must have two arguments")
	}

	var newPoints int
	var userId = args[0]

	if function == "addPoints" {

		newPoints = t.addPoints(stub, args)

	}

	if function == "deletePonts" {
		newPoints = t.deletePoints(stub, args)
	}

	fmt.Println(newPoints)

	userData, err := stub.GetState(userId)

	if err != nil {
		return nil, errors.New("Exiting from Invoke with error. stub.GetState gave some problem")
	}

	fmt.Println(userData)

	var userRetrieved User
	err = json.Unmarshal(userData, &userRetrieved)
	userRetrieved.totalPoint = newPoints

	var userRetrievedAddress = &userRetrieved

	userDataBytes, err = json.Marshal(userRetrievedAddress)

	err = stub.PutState(userId, userDataBytes)

	return nil, nil

}

func (t *SimpleChaincode) addPoints(stub shim.ChaincodeStubInterface, args []string) int {

	var points = t.getPoints(stub, args)
	var numberOfPointsToAdd = args[1]
	numberOfPointsInt, err := strconv.Atoi(numberOfPointsToAdd)
	if err != nil {
		return 0
	}
	points = points + numberOfPointsInt
	fmt.Println(points)

	return points

}

func (t *SimpleChaincode) deletePoints(stub shim.ChaincodeStubInterface, args []string) int {

	var points = t.getPoints(stub, args)

	var numberOfPointsToSub = args[1]
	numberOfPointsInt, err := strconv.Atoi(numberOfPointsToSub)
	if err != nil {
		return 0
	}

	points = points - numberOfPointsInt

	fmt.Println(points)

	return points

}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var currentPoints int
	if len(args) != 1 {
		return nil, errors.New("must have only 1 arguments")
	}

	if function == "getPoints" {
		currentPoints = t.getPoints(stub, args)
	}
	if currentPoints < 0 {
		return nil, errors.New("returning from query with error. current points are negative. Add points")
	}
	return nil, nil
}

func (t *SimpleChaincode) getPoints(stub shim.ChaincodeStubInterface, args []string) int {

	var userId = args[0]

	fmt.Println(userId)

	bytes, err := stub.GetState(userId)

	if err != nil {
		return 0
	}
	var userRetrieved User
	err = json.Unmarshal(bytes, &userRetrieved)
	var currentPoints = userRetrieved.totalPoint

	fmt.Println(currentPoints)
	return currentPoints

}
