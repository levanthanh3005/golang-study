package main
import (
	"fmt"
	"encoding/json"
	// "strconv"
	contract "myapp/ucontract/fabcar"
	contractapi "myapp/utility"
	)

func example1() {
    fmt.Println("Simple test")

	// var chaincode *SmartContract
	// chaincode = new(SmartContract)
	// chaincode.Name = "aa"
	// chaincode.Test = 1
	// fmt.Println(chaincode.Name)
    sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)

    // stub.mapData = map[string][]byte
	// b := []byte(`{"test":"1"}`)
	// var f interface{}
	// err := json.Unmarshal(b, &f)
	
	// stub.mapData = f.(map[string]interface{})
	fmt.Println(stub)

    ctx.SetStub(stub)

    carAsBytes, _ := json.Marshal(`{"test":"1"}`)
    stub.InitState()
	stub.PutState("CAR1", carAsBytes)

    // fmt.Println(ctx.GetStub())
    // ctx.stub.
    sm.InitLedger(ctx)

    fmt.Println(ctx.GetStub().GetAllState())

    b,_ := ctx.GetStub().GetState("CAR0")
	var f interface{}
	json.Unmarshal(b, &f)

	fmt.Println(f.(map[string]interface{}))
}

func example2() {
    fmt.Println("CreateCar => GetState and QueryCar")

    sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)
    stub.InitState()
    ctx.SetStub(stub)
    sm.InitLedger(ctx)
    sm.CreateCar(ctx, "carNumber", "make", "model", "colour", "owner")

    fmt.Println("GetState")
    b,_ := ctx.GetStub().GetState("carNumber")
	var f interface{}
	json.Unmarshal(b, &f)

	fmt.Println(f.(map[string]interface{}))

    fmt.Println("QueryCar")
	rs,_ := sm.QueryCar(ctx, "carNumber")
    inrec, _ := json.Marshal(rs)
	json.Unmarshal(inrec, &f)
	fmt.Println(f.(map[string]interface{}))
}

func example3() {
    fmt.Println("CreateCar => GetState and RemoveCar => QueryCar")

    sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)
    stub.InitState()
    ctx.SetStub(stub)
    sm.InitLedger(ctx)
    sm.CreateCar(ctx, "carNumber", "make", "model", "colour", "owner")

    fmt.Println("GetState")
    b,_ := ctx.GetStub().GetState("carNumber")
	var f interface{}
	json.Unmarshal(b, &f)

	fmt.Println(f.(map[string]interface{}))

    fmt.Println("RemoveCar")
	err := sm.RemoveCar(ctx, "carNumber")
	fmt.Println(err)

    fmt.Println("QueryCar")
	rs,err := sm.QueryCar(ctx, "carNumber")
	fmt.Println(err)
	fmt.Println(rs)
 //    inrec, _ := json.Marshal(rs)
	// json.Unmarshal(inrec, &f)
	// fmt.Println(f.(map[string]interface{}))
}

func example4() {
    fmt.Println("CreatePrivateData => Query")

    sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)
    stub.InitState()
    ctx.SetStub(stub)
    sm.InitLedger(ctx)
    sm.CreateCar(ctx, "carNumber", "make", "model", "colour", "owner")

    fmt.Println("GetState")
    b,_ := ctx.GetStub().GetState("carNumber")
	var f interface{}
	json.Unmarshal(b, &f)
	fmt.Println(f.(map[string]interface{}))

	sm.CreatePrivateDataForCar(ctx,"test","carNumber","abc")
	rs,err := sm.ReadPrivateDataForCar(ctx, "test", "carNumber")
	fmt.Println(err)
	fmt.Println(rs)
}

func example5() {
    fmt.Println("Client Identity")
    sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)
    stub.InitState()
    ctx.SetStub(stub)
    
    sm.InitLedger(ctx)
    sm.CreateCar(ctx, "carNumber", "make", "model", "colour", "owner")

    fmt.Println(sm.GetSenderId(ctx))
    // fmt.Println()
}

func main() {
	// example1()
	example2()
	example3()
	example4()
	example5()
}