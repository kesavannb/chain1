package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

var marbleIndexStr = "_marbleindex"				//name for the key/value that will store a list of all known marbles


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(marbleIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "add_name" {									//create a new marble
		return t.add_name(stub, args)
	}
	
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) add_name(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	name := args[0]

	//get the name index
	marblesAsBytes, err := stub.GetState(marbleIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get name index")
	}
	var marbleIndex []string
	json.Unmarshal(marblesAsBytes, &marbleIndex)							//un stringify it aka JSON.parse()
	
	//append
	marbleIndex = append(marbleIndex, name)									//add marble name to index list
	fmt.Println("! marble index: ", marbleIndex)
	jsonAsBytes, _ := json.Marshal(marbleIndex)
	err = stub.PutState(marbleIndexStr, jsonAsBytes)						//store name of marble

	fmt.Println("- end init marble")
	return nil, nil
}
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	name := args[0]

	//get the name index
	marblesAsBytes, err := stub.GetState(marbleIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get name index")
	}
	var marbleIndex []string
	json.Unmarshal(marblesAsBytes, &marbleIndex)							//un stringify it aka JSON.parse()
	
	//append
	marbleIndex = append(marbleIndex, name)									//add marble name to index list

		fmt.Println("marbleIndex",marbleIndex)
	
	return nil, nil
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
