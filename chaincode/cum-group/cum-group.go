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

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Group structure, with several properties.
Structure tags are used by encoding/json library
*/
type Group struct {
	GroupId   string `json:"groupId"`
	GroupName string `json:"groupName"`
	GroupDesc string `json:"groupDesc"`
}

/*
 * The Init method *
 called when the Smart Contract "elza-chaincode" is instantiated by the network
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
	if function == "queryGroupById" {
		return s.queryGroupById(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "addGroup" {
		return s.addGroup(APIstub, args)
	} else if function == "queryAllGroups" {
		return s.queryAllGroups(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryParsel method *
 Used to view the records of one particular parsel
 It takes one argument -- the key for the parsel in question
*/
func (s *SmartContract) queryGroupById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	groupAsBytes, _ := APIstub.GetState(args[0])
	if groupAsBytes == nil {
		return shim.Error("Could not locate parsel")
	}
	return shim.Success(groupAsBytes)
}

/*
 * The initLedger method *
Will add test data to our network
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	// groups := []Group{
	// 	Group{GroupId: "001", GroupName: "AC17", GroupDesc: "Description for AB17"},
	// 	Group{GroupId: "002", GroupName: "AC18", GroupDesc: "Description for AB18"},
	// }

	// i := 0
	// for i < len(groups) {
	// 	fmt.Println("i is ", i)
	// 	groupAsBytes, _ := json.Marshal(groups[i])
	// 	APIstub.PutState(fmt.Sprintf("%X", rand.Int()), groupAsBytes)
	// 	fmt.Println("Added", groups[i])
	// 	i = i + 1
	// }

	return shim.Success(nil)
}

/*
 * The addGroup method
 * This method takes in four arguments (attributes to be saved in the ledger).
 */
func (s *SmartContract) addGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var group = Group{GroupId: args[1], GroupName: args[2], GroupDesc: args[3]}

	groupAsBytes, _ := json.Marshal(group)
	err := APIstub.PutState(fmt.Sprintf("%X", rand.Int()), groupAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record new group: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllGroups method *
allows for assessing all the records added to the ledger(all groups in the delivery system)
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

	fmt.Printf("- queryAllGroups:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 *   main function : calls the Start function
 *   The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
