// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
*/

package main

/* Imports
* 5 utility libraries for handling bytes, reading and writing JSON,
formatting, and string manipulation
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Student structure, with several properties.
Structure tags are used by encoding/json library
*/
type Student struct {
	StudentId   string `json:"studentId"`
	StudentName string `json:"studentName"`
	GroupName   string `json:"groupName"`
}

/* Define Student Test structure, with several properties.
Structure tags are used by encoding/json library
*/
type Stest struct {
	StestId string `json:"testId"`

	Group string `json:"group"`

	Course string `json:"course"`

	Teacher string `json:"teacher"`

	StudentId string `json:"studentId"`

	Student string `json:"student"`

	Rate string `json:"rate"`

	StestTS string `json:"stestTS"`

	StestDesc string `json:"stestDesc"`
}

/*
 * The Init method *
 called when the Smart Contract "elza-rec" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function
 -- see initLedger()
*/
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "posta-chaincode"
 The app also specifies the specific smart contract function to call with args
*/
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAllTests" {
		return s.queryAllTests(APIstub)
	} else if function == "queryTestById" {
		return s.queryTestById(APIstub, args)
	} else if function == "createTestForGroup" {
		return s.createTestForGroup(APIstub, args)
	} else if function == "queryTestByStudent" {
		return s.queryTestByStudent(APIstub, args)
	} else if function == "prepareForExam" {
		return s.prepareForExam(APIstub, args)
	} else if function == "takeTheTest" {
		return s.takeTheTest(APIstub, args)
	} else if function == "addStudent" {
		return s.addStudent(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The initLedger method *
Will add test data to our network
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	stests := []Stest{
		Stest{StestId: "001", Group: "AB17", Course: "Maths", Teacher: "Ivanov", Student: "AB1701", Rate: "", StestTS: "", StestDesc: ""},
		Stest{StestId: "002", Group: "AB17", Course: "Phisycs", Teacher: "Petrov", Student: "AB1701", Rate: "", StestTS: "", StestDesc: ""},
	}

	i := 0
	for i < len(stests) {
		fmt.Println("i is ", i)
		stestAsBytes, _ := json.Marshal(stests[i])
		APIstub.PutState(fmt.Sprintf("%X", rand.Int()), stestAsBytes)
		fmt.Println("Added", stests[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/* The addStudent method *
   Generate initial student record
   This method takes in four arguments (attributes to be saved in the ledger).
*/
func (s *SmartContract) addStudent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Starting index for records

	var student = Student{StudentId: args[0], StudentName: args[1], GroupName: args[2]}

	studentAsBytes, _ := json.Marshal(student)

	err := APIstub.PutState(fmt.Sprintf("%X", rand.Int()), studentAsBytes)

	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to add student to system with id: %s", args[0]))
	}

	fmt.Println("Added sudent with id: %s ", args[0])

	return shim.Success(nil)
}

/*
 * The queryAllTests method *
allows for assessing all the records added to the ledger(all groups in the delivery system)
This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SmartContract) queryAllTests(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "9999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTests:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The queryTestById method *
 Used to view the records of one particular parsel
 It takes one argument -- the key for the parsel in question
*/

func (s *SmartContract) queryTestById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	stestAsBytes, _ := APIstub.GetState(args[0])
	if stestAsBytes == nil {
		return shim.Error("Could not locate test record")
	}
	return shim.Success(stestAsBytes)
}

/* The createTestForGroup method *
   Generate list of records for one exam/course/teacher
   This method takes in four arguments (attributes to be saved in the ledger).
*/
func (s *SmartContract) createTestForGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	i := 0

	// Counter of added records
	var counter int
	counter, _ = strconv.Atoi(args[2])

	// Starting index for records
	var startIndex int
	startIndex, _ = strconv.Atoi(args[0])

	for i < counter {

		// Combine student id
		studentName := args[1] + strconv.Itoa(i)

		var stest = Stest{StestId: strconv.Itoa(startIndex + i), Group: args[1], Course: args[3], Teacher: args[4], Student: studentName, Rate: "", StestTS: "", StestDesc: ""}

		stestAsBytes, _ := json.Marshal(stest)

		err := APIstub.PutState(fmt.Sprintf("%X", rand.Int()), stestAsBytes)

		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to record new test for the group: %s", args[0]))
		}

		fmt.Println("Added", stest)
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The queryTestByStudent method *
   allows for assessing all the records from selected student

    Returns JSON string containing results.
*/

func (s *SmartContract) queryTestByStudent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	startKey := "0"
	endKey := "9999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Create an object
		stest := Stest{}
		// Unmarshal record to stest object
		json.Unmarshal(queryResponse.Value, &stest)

		// Add only filtered ny sender records
		if stest.Student == args[0] {

			// Add comma before array members,suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}

			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}
	}
	buffer.WriteString("]")

	if bArrayMemberAlreadyWritten == false {
		return shim.Error("No tests for student")
	}

	fmt.Printf("- queryTestByStudent:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The prepareForExam method *
   allows for assessing all the records from selected group/course

	 Returns JSON string containing results.
*/

func (s *SmartContract) prepareForExam(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	startKey := "0"
	endKey := "9999"

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Create an object
		stest := Stest{}
		// Unmarshal record to stest object
		json.Unmarshal(queryResponse.Value, &stest)

		// Add only filtered ny sender records

		if stest.Group == args[0] && stest.Course == args[1] {

			// Add comma before array members,suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}

			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}
	}
	buffer.WriteString("]")

	if bArrayMemberAlreadyWritten == false {
		return shim.Error("No group/course found")
	}

	fmt.Printf("- prepareForExam:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The takeTheTest method *
 * The data in the stest state can be updated .
 */
func (s *SmartContract) takeTheTest(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	stestAsBytes, _ := APIstub.GetState(args[0])

	if stestAsBytes == nil {
		return shim.Error("Could not locate unpassed test")
	}

	stest := Stest{}

	json.Unmarshal(stestAsBytes, &stest)

	if stest.Rate != "" {
		return shim.Error("Could not locate unpassed test")
	}

	// Normally check that the specified argument is a valid participant of exam
	// we are skipping this check for this example
	stest.StestTS = time.Now().Format(time.RFC1123Z)
	stest.Rate = args[3]

	stestAsBytes, _ = json.Marshal(stest)
	err := APIstub.PutState(args[0], stestAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change status of exam: %s", args[0]))
	}

	return shim.Success(stestAsBytes)
}

/*
 * main function  - calls the Start function
   The main function starts the chaincode in the container during instantiation.
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
