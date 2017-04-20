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
	ffId        int    `json:"ffid"`
	title       string `json:"title"`
	gender      string `json:"gender"`
	firstName   string `json:"firstname"`
	lastName    string `json:"lastname"`
	dob         string `json:"dob"`
	email       string `json:"email"`
	createdBy   string `json:"createdBy"`
	totalPoint  int    `json:"totalPoint"`
	country     string `json:"country"`
	addressLine string `json:"addressLine"`
	city        string `json:"city"`
	zip         string `json:"zip"`
}

/*func (t *SimpleChaincode) registerNewUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("inside addNew User")

	var newUser User
	var newUserByteArray []byte
	var err error
	var id, tPoint int
	var newUserDataAddress *User

	if len(args) != 13 {
		return nil, errors.New("must have thirteen arguments")
	}

		var idString string
		var tpoint string

		idString = args[0]
		fmt.Println("ID String")
		fmt.Print(idString)

		tpoint = args[12]
		fmt.Println("tPoint")
		fmt.Print(tpoint)

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

		newUserDataAddress = &newUser

	newUserByteArray, err = json.Marshal(newUserDataAddress)

	err = stub.PutState(idString, newUserByteArray)

	if err != nil {
		return nil, errors.New("data cannot be pushed successfully. returned from registerNewUser")
	}

	return newUserByteArray, nil

}
*/

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	//var newUserByteArray []byte
	var userId string
	var userDetails string

	fmt.Println("Inside Init . This is used to create an new User")
	fmt.Println("Inside registerNewUser")

	userId = args[0]
	userDetails = args[1]

	err = stub.PutState(userId, []byte(userDetails))

	if err != nil {
		fmt.Println("Could not save userDetails to ledger", err)
		return nil, err

	}
	//newUserByteArray, err = t.registerNewUser(stub, args)

	return nil, nil

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	var newPoints int
	var userDataBytes []byte
	var points int

	if len(args) != 2 {
		return nil, errors.New("must have two arguments")
	}

	var userId = args[0]

	fmt.Println("userId : ")
	fmt.Print(userId)

	var existingUser User

	userDataBytes, err = stub.GetState(userId)

	if err != nil {
		return nil, errors.New("Error retreiving User with given Id")
	}

	err = json.Unmarshal(userDataBytes, &existingUser)

	if err != nil {
		return nil, errors.New("Problem in unmarshalling data")
	}

	fmt.Print("existingUser.totalPoint ")
	fmt.Println(existingUser.totalPoint)

	points = existingUser.totalPoint

	if function == "addPoints" {

		fmt.Println("Inside addPoints")

		newPoints = t.addPoints(stub, args, points)

	}

	if function == "deletePonts" {

		fmt.Println("Inside deletePoints")

		newPoints = t.deletePoints(stub, args, points)

	}
	fmt.Println("back from addPoints method inside invoke . points to be added are : ")

	existingUser.totalPoint = newPoints
	userDataBytes, err = json.Marshal(&existingUser)

	err = stub.PutState(userId, userDataBytes)

	return userDataBytes, nil

}

func (t *SimpleChaincode) addPoints(stub shim.ChaincodeStubInterface, args []string, points int) int {

	var numberOfPointsToAdd = args[1]
	numberOfPointsInt, err := strconv.Atoi(numberOfPointsToAdd)

	if err != nil {

		return 0
	}

	fmt.Print("number of Points to add")
	fmt.Println(numberOfPointsToAdd)

	fmt.Print("number of Points currently")
	fmt.Println(points)

	points = points + numberOfPointsInt

	fmt.Println("Points added and final result is ")
	fmt.Print(points)

	return points

}

func (t *SimpleChaincode) deletePoints(stub shim.ChaincodeStubInterface, args []string, points int) int {

	var numberOfPointsToSub = args[1]
	numberOfPointsInt, err := strconv.Atoi(numberOfPointsToSub)

	if err != nil {

		return 0
	}

	fmt.Print("number of Points to subtract")
	fmt.Println(numberOfPointsToSub)

	fmt.Print("number of Points Int currently")
	fmt.Println(points)

	points = points - numberOfPointsInt

	fmt.Println("Points subtracted and final result is ")
	fmt.Print(points)

	return points

}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var currentPoints []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("must have only 1 arguments")
	}

	if function == "getPoints" {
		currentPoints, err = t.getPoints(stub, args)
	}
	if err != nil {
		return nil, errors.New("returning from query with error. current points are not correct. Add points")
	}
	return currentPoints, nil
}

func (t *SimpleChaincode) getPoints(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var userId = args[0]

	fmt.Println(userId)

	bytes, err := stub.GetState(userId)

	if err != nil {
		var userRetrieved User
		err = json.Unmarshal(bytes, &userRetrieved)
		var currentPoints = userRetrieved.totalPoint

		fmt.Println("currentPoints ")
		fmt.Print(currentPoints)
	}
	return bytes, nil

}
