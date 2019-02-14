/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the invoice structure, with 10 properties.  Structure tags are used by encoding/json library
type Invoice struct {
	InvoiceNumber   string  `json:"invoicenumber"`
	BilledTo        string  `json:"billedto"`
	InvoiceDate     string  `json:"invoicedate"`
	InvoiceAmount   float64 `json:"invoiceamount"`
	ItemDescription string  `json:"itemdescription"`
	GR              bool    `json:"gr"`
	IsPaid          bool    `json:"ispaid"`
	PaidAmount      float64 `json:"paidamount"`
	IsRepaid        bool    `json:"isrepaid"`
	RepaymentAmount float64 `json:"repaymentamount"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createInvoice" {
		return s.createInvoice(APIstub, args)
	} else if function == "displayAllInvoice" {
		return s.displayAllInvoice(APIstub)
	} else if function == "isGoodReceived" {
		return s.isGoodReceived(APIstub, args)
	} else if function == "isPaid" {
		return s.isPaid(APIstub, args)
	} else if function == "isRepaid" {
		return s.isRepaid(APIstub, args)
	} else if function == "getAuditHistoryForInvoice" {
		return s.getAuditHistoryForInvoice(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// function for initializing the ledger
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	// initialize invoice
	invoice := []Invoice{
		Invoice{
			InvoiceNumber:   "1001",
			BilledTo:        "IBM",
			InvoiceDate:     "10FEB2019",
			InvoiceAmount:   1000.00,
			ItemDescription: "Phone",
			GR:              false,
			IsPaid:          false,
			PaidAmount:      0.00,
			IsRepaid:        false,
			RepaymentAmount: 0.00,
		},
	}

	i := 0
	for i < len(invoice) {
		fmt.Println("i is ", i)
		invoiceAsBytes, _ := json.Marshal(invoice[i])
		APIstub.PutState("INVOICE"+strconv.Itoa(i), invoiceAsBytes)
		fmt.Println("Added", invoice[i])
		i = i + 1
	}

	return shim.Success(nil)
}

// function for creating invoice
func (s *SmartContract) createInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	invoiceAmount, _ := strconv.ParseFloat(args[4], 64)
	// initialize all arguments

	var invoice = Invoice{
		InvoiceNumber:   args[1],
		BilledTo:        args[2],
		InvoiceDate:     args[3],
		InvoiceAmount:   invoiceAmount,
		ItemDescription: args[5],

		GR:              false,
		IsPaid:          false,
		PaidAmount:      0.00,
		IsRepaid:        false,
		RepaymentAmount: 0.00,
	}

	invoiceAsBytes, _ := json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(invoiceAsBytes)
}

// function for displaying all invoice
func (s *SmartContract) displayAllInvoice(APIstub shim.ChaincodeStubInterface) sc.Response {

	// initialize range of results
	startKey := "INVOICE0"
	endKey := "INVOICE999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	// display results matched
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
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

	fmt.Printf("- displayAllInvoice:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// function for checking goods received
func (s *SmartContract) isGoodReceived(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// invoice id and goods received
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.GR = true

	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}

// function for checking if paid to supplier
func (s *SmartContract) isPaid(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// invoiceid and ispaid
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	paidAmount, _ := strconv.ParseFloat(args[1], 64)
	json.Unmarshal(invoiceAsBytes, &invoice)

	if paidAmount < invoice.InvoiceAmount {
		invoice.PaidAmount = paidAmount
		invoice.IsPaid = true
	} else {
		return shim.Error("Paid amount must be always less than the invoice amount.")
	}

	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}

// function for checking if paid to bank
func (s *SmartContract) isRepaid(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// invoiceid and isrepaid
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	repaymentAmount, _ := strconv.ParseFloat(args[1], 64)
	json.Unmarshal(invoiceAsBytes, &invoice)

	if repaymentAmount < invoice.RepaymentAmount {
		invoice.RepaymentAmount = repaymentAmount
		invoice.IsRepaid = true
	} else {
		return shim.Error("Paid amount must be always less than the invoice amount.")
	}

	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}

// function for getting invoice audit history
func (s *SmartContract) getAuditHistoryForInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	invoiceKey := args[0]

	resultsIterator, err := APIstub.GetHistoryForKey(invoiceKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the invoice
	var buffer bytes.Buffer
	buffer.WriteString("[")

	// display audit history if results found
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		buffer.WriteString(string(response.Value))

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

