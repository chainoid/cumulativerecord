// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
*/

package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
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
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "getStudentRecord" {
		return s.getStudentRecord(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The initLedger method
 *
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

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
