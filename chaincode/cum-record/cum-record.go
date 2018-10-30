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
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define Enum for classificate record type
type RecordType string

const (
	GroupType   RecordType = "G"
	StudentType RecordType = "S"
	TeacherType RecordType = "T" // For future use
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Student structure, with several properties.
Structure tags are used by encoding/json library
*/
type StudentRecord struct {
	RecordType  string        `json:"recordType"`
	GroupName   string        `json:"groupName"`
	StudentId   string        `json:"studentId"`
	StudentName string        `json:"studentName"`
	Description string        `json:"description"`
	RegisterTS  string        `json:"registerTS"`
	RecordList  []StudentTest `json:"recordList"`
}

/* Define Student structure, with several properties.
Structure tags are used by encoding/json library
*/
type Student struct {
	StudentId   string `json:"studentId"`
	StudentName string `json:"studentName"`
	GroupName   string `json:"groupName"`
}

/* Define Student Test structure, with several properties
Structure tags are used by encoding/json library
*/
type StudentTest struct {
	StestId     string `json:"testId"`
	Group       string `json:"group"`
	Course      string `json:"course"`
	Teacher     string `json:"teacher"`
	AssignedTS  string `json:"assignedTS"`
	Rate        string `json:"rate"`
	ExecuteTS   string `json:"executeTS"`
	ExecuteDesc string `json:"executeDesc"`
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
	} else if function == "queryAllGroups" {
		return s.queryAllGroups(APIstub)
	} else if function == "addGroup" {
		return s.addGroup(APIstub, args)
	} else if function == "queryAllStudents" {
		return s.queryAllStudents(APIstub)
	} else if function == "queryAllTests" {
		return s.queryAllTests(APIstub)
	} else if function == "queryTestById" {
		return s.queryTestById(APIstub, args)
	} else if function == "createTestForGroup" {
		return s.createTestForGroup(APIstub, args)
	} else if function == "getStudentRecord" {
		return s.getStudentRecord(APIstub, args)
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
 * The initLedger method
 * Will add group and student records to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	stests := []StudentRecord{
		StudentRecord{RecordType: "G", GroupName: "AB17", Description: "Desc AB17"},
		StudentRecord{RecordType: "G", GroupName: "AB18", Description: "Desc AB18"},
		StudentRecord{RecordType: "G", GroupName: "AB19", Description: "Desc AB19"},
		StudentRecord{RecordType: "G", GroupName: "AB20", Description: "Desc AB20"},

		StudentRecord{RecordType: "S", GroupName: "AB17", StudentId: "AB1701", StudentName: "Student 1701", RegisterTS: time.Now().Format(time.RFC3339), Description: "Desc 1701"},
		StudentRecord{RecordType: "S", GroupName: "AB17", StudentId: "AB1702", StudentName: "Student 1702", RegisterTS: time.Now().Format(time.RFC3339), Description: "Desc 1702"},
		StudentRecord{RecordType: "S", GroupName: "AB17", StudentId: "AB1703", StudentName: "Student 1703", RegisterTS: time.Now().Format(time.RFC3339), Description: "Desc 1703"},
		StudentRecord{RecordType: "S", GroupName: "AB17", StudentId: "AB1704", StudentName: "Student 1704", RegisterTS: time.Now().Format(time.RFC3339), Description: "Desc 1704"},
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

/*
 * The queryAllGroups method *
allows for assessing all group records added to the ledger(all groups in the cumulative system)
This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SmartContract) queryAllGroups(APIstub shim.ChaincodeStubInterface) sc.Response {

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
		studentRecord := StudentRecord{}
		// Unmarshal record to stest object
		json.Unmarshal(queryResponse.Value, &studentRecord)

		// Add only filtered by RecordType as Group records
		if studentRecord.RecordType == "G" {

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

	fmt.Printf("- queryAllGroups:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The addGroup method
 * This method takes in four arguments (attributes to be saved in the ledger).
 */
func (s *SmartContract) addGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var groupRecord = StudentRecord{RecordType: "G", GroupName: args[0], Description: args[1]}

	groupRecordAsBytes, _ := json.Marshal(groupRecord)
	err := APIstub.PutState(fmt.Sprintf("%X", rand.Int()), groupRecordAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record new group: %s", args[0]))
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

	var studentRecord = StudentRecord{RecordType: "S", StudentId: args[0], StudentName: args[1], GroupName: args[2], Description: args[3], RegisterTS: time.Now().Format(time.RFC3339)}

	studentRecordAsBytes, _ := json.Marshal(studentRecord)

	err := APIstub.PutState(fmt.Sprintf("%X", rand.Int()), studentRecordAsBytes)

	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to add student to system with id: %s", args[0]))
	}

	fmt.Println("Added student with id: ", args[0])

	return shim.Success(nil)
}

/*
 * The queryAllStudents method *
allows for assessing all student records added to the ledger
This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SmartContract) queryAllStudents(APIstub shim.ChaincodeStubInterface) sc.Response {

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
		studentRecord := StudentRecord{}
		// Unmarshal record to stest object
		json.Unmarshal(queryResponse.Value, &studentRecord)

		// Add only filtered by RecordType as Group records
		if studentRecord.RecordType == "S" {

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

	fmt.Printf("- queryAllStudents:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The queryAllTests method *
allows for assessing all the records added to the ledger(all groups in the delivery system)
This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SmartContract) queryAllTests(APIstub shim.ChaincodeStubInterface) sc.Response {

	// startKey := "0"
	// endKey := "9999"

	// resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// defer resultsIterator.Close()

	// // buffer is a JSON array containing QueryResults
	// var buffer bytes.Buffer
	// buffer.WriteString("[")

	// bArrayMemberAlreadyWritten := false
	// for resultsIterator.HasNext() {
	// 	queryResponse, err := resultsIterator.Next()
	// 	if err != nil {
	// 		return shim.Error(err.Error())
	// 	}
	// 	// Add comma before array members,suppress it for the first array member
	// 	if bArrayMemberAlreadyWritten == true {
	// 		buffer.WriteString(",")
	// 	}
	// 	buffer.WriteString("{\"Key\":")
	// 	buffer.WriteString("\"")
	// 	buffer.WriteString(queryResponse.Key)
	// 	buffer.WriteString("\"")

	// 	buffer.WriteString(", \"Record\":")
	// 	// Record is a JSON object, so we write as-is
	// 	buffer.WriteString(string(queryResponse.Value))
	// 	buffer.WriteString("}")
	// 	bArrayMemberAlreadyWritten = true
	// }
	// buffer.WriteString("]")

	// fmt.Printf("- queryAllTests:\n%s\n", buffer.String())

	// return shim.Success(buffer.Bytes())

	// TODO remove
	return shim.Success(nil)
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
   Generate list of records for one group/course/teacher
   This method takes in four arguments (attributes to be saved in the ledger).
*/
func (s *SmartContract) createTestForGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	startKey := "0"
	endKey := "9999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Create an object
		studentRecord := StudentRecord{}
		// Unmarshal record to stest object
		json.Unmarshal(queryResponse.Value, &studentRecord)

		if studentRecord.GroupName == args[0] {
			var studentTest = StudentTest{StestId: fmt.Sprintf("%X", rand.Int()), Group: args[0], Course: args[1], Teacher: args[2],
				AssignedTS: time.Now().Format(time.RFC1123Z), ExecuteDesc: ""}

			studentRecord.RecordList = append(studentRecord.RecordList, studentTest)

			studentRecordAsBytes, _ := json.Marshal(studentRecord)

			APIstub.PutState(queryResponse.Key, studentRecordAsBytes)

			fmt.Println("Added", studentTest)
		}

	}

	return shim.Success(nil)
}

/*
 * The getStudentRecord method *
   allows for assessing all the records from selected student

    Returns JSON string containing results.
*/

func (s *SmartContract) getStudentRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	studentRecordAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Could not locate student data")
	}

	return shim.Success(studentRecordAsBytes)
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

	fmt.Printf("- prepareForExam param args[0] :%s  args[1] :%s\n", args[0], args[1])

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing Results
	var buffer bytes.Buffer

	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	// Iteration my the Student Record List
	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Create an object
		studentRecord := StudentRecord{}
		// Unmarshal record to stest object
		json.Unmarshal(queryResponse.Value, &studentRecord)

		// Selection by the Record type and Group name
		if studentRecord.RecordType == "S" && studentRecord.GroupName == args[0] {

			// Iteration by the Course

			for i := 0; i < len(studentRecord.RecordList); i++ {

				if studentRecord.RecordList[i].Course == args[1] {

					// Add comma before array members,suppress it for the first array member
					if bArrayMemberAlreadyWritten == true {
						buffer.WriteString(",")
					}

					buffer.WriteString("{\"Key\":")
					buffer.WriteString("\"")
					buffer.WriteString(queryResponse.Key)
					buffer.WriteString("\"")
					buffer.WriteString(", \"Record\":")

					// Put only selected fields
					buffer.WriteString("{\"studentId\":\"")
					buffer.WriteString(queryResponse.Key)
					buffer.WriteString("\",")

					buffer.WriteString("\"testId\":\"")
					buffer.WriteString(studentRecord.RecordList[i].StestId)
					buffer.WriteString("\",")

					buffer.WriteString("\"studentName\":\"")
					buffer.WriteString(studentRecord.StudentName)
					buffer.WriteString("\",")

					buffer.WriteString("\"group\":\"")
					buffer.WriteString(studentRecord.GroupName)
					buffer.WriteString("\",")

					buffer.WriteString("\"course\":\"")
					buffer.WriteString(studentRecord.RecordList[i].Course)
					buffer.WriteString("\",")

					buffer.WriteString("\"assignedTS\":\"")
					buffer.WriteString(studentRecord.RecordList[i].AssignedTS)
					buffer.WriteString("\",")

					buffer.WriteString("\"teacher\":\"")
					buffer.WriteString(studentRecord.RecordList[i].Teacher)
					buffer.WriteString("\",")

					buffer.WriteString("\"executeTS\":\"")
					buffer.WriteString(studentRecord.RecordList[i].ExecuteTS)
					buffer.WriteString("\",")

					buffer.WriteString("\"rate\":\"")
					buffer.WriteString(studentRecord.RecordList[i].Rate)
					buffer.WriteString("\"")

					buffer.WriteString("}")

					buffer.WriteString("}")
					bArrayMemberAlreadyWritten = true
				}
			}
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

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	fmt.Printf("- takeTheTest param args[0] :%s  args[1] :%s args[3] :%s\n", args[0], args[1], args[2])

	studentRecordAsBytes, _ := APIstub.GetState(args[0])

	if studentRecordAsBytes == nil {
		return shim.Error("Could not locate selected student record")
	}

	studentRecord := StudentRecord{}

	json.Unmarshal(studentRecordAsBytes, &studentRecord)

	// Looking for selected test by course name
RecordListLoop:

	for i := 0; i < len(studentRecord.RecordList); i++ {

		if studentRecord.RecordList[i].Course == args[1] {

			if studentRecord.RecordList[i].Rate != "" {
				return shim.Error("Selected test already passed.")
			}

			studentRecord.RecordList[i].Rate = args[2]
			studentRecord.RecordList[i].ExecuteTS = time.Now().Format(time.RFC1123Z)
			//studentRecord.RecordList[i].ExecuteDesc = args[3]

			studentRecordAsBytes, _ := json.Marshal(studentRecord)

			err := APIstub.PutState(args[0], studentRecordAsBytes)

			if err != nil {
				return shim.Error(fmt.Sprintf("Failed to change status of student record %s", args[0]))
			}

			break RecordListLoop
		}
	}

	return shim.Success(studentRecordAsBytes)
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
