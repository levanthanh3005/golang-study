package main
import (
	"fmt"
	"myapp/functions"
	)

// SmartContract provides functions for managing a car
type SmartContract struct {
	functions.Contract
}

func main() {
    fmt.Println("hello world")
    fmt.Println(functions.GetValue(1))

	// var chaincode *SmartContract
	// chaincode = new(SmartContract)
	// chaincode.Name = "aa"
	// chaincode.Test = 1
	// fmt.Println(chaincode.Name)

	chaincode, err := functions.NewChaincode(new(SmartContract))

	fmt.Println(chaincode)
	fmt.Println(err)

}