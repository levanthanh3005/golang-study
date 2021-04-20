package main
import (
	"fmt"
	// "encoding/json"
	// "strconv"
	contract "myapp/ucontract/ransharing-private"
	contractapi "myapp/utility"
	)

func example1() {
	fmt.Println("Create Cells and check exception of duplication of Cell")
	sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)
    stub.InitState()
    ctx.SetStub(stub)
    
    sm.InitLedger(ctx)
    fmt.Println(sm.GetSenderId(ctx))
    err := sm.CreateCellData(ctx, "test",`{"CellId" : "cell-id1"}`,`{"CellId" : "cell-id1","RAT" : "1111","DL_EARFCN" : "1111","PCI" : "111"}`)
    fmt.Println(err)
    ls,err := sm.GetCellSet(ctx,"test")
    fmt.Println(ls)
    fmt.Println(err)

    err = sm.CreateCellData(ctx, "test",`{"CellId" : "cell-id2"}`,`{"CellId" : "cell-id2","RAT" : "2222","DL_EARFCN" : "2222","PCI" : "222"}`)
    fmt.Println(err)
    ls,err = sm.GetCellSet(ctx,"test")
    fmt.Println(ls)
    fmt.Println(err)

    fmt.Println("###########")

    err = sm.CreateCellData(ctx, "test",`{"CellId" : "cell-id1"}`,`{"CellId" : "cell-id1","RAT" : "1111","DL_EARFCN" : "1111","PCI" : "111"}`)
    fmt.Println(err)

    fmt.Println("###########")

    ls,err = sm.GetCellSet(ctx,"test")
    fmt.Println(ls)
    fmt.Println(err)
}

func example2() {
	fmt.Println("Create Cells and Update")
	sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)
    stub.InitState()
    ctx.SetStub(stub)
    
    sm.InitLedger(ctx)
    fmt.Println(sm.GetSenderId(ctx))
    err := sm.CreateCellData(ctx, "test",`{"CellId" : "cell-id1"}`,`{"CellId" : "cell-id1","RAT" : "1111","DL_EARFCN" : "1111","PCI" : "111"}`)
    fmt.Println(err)
    ls,err := sm.GetCellSet(ctx,"test")
    fmt.Println(ls)
    fmt.Println(err)

    err = sm.UpdateCellData(ctx, "test",`{"CellId" : "cell-id1"}`,`{"CellId" : "cell-id1","RAT" : "2222","DL_EARFCN" : "2222","PCI" : "222"}`)
    fmt.Println(err)
    ls,err = sm.GetCellSet(ctx,"test")
    fmt.Println(ls)
    fmt.Println(err)
}

func example3() {
	fmt.Println("Create Cell Neigbour Relation")
	sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)
    stub.InitState()
    ctx.SetStub(stub)
    
    sm.InitLedger(ctx)
    fmt.Println(sm.GetSenderId(ctx))
    err := sm.CreateCellData(ctx, "test",`{"CellId" : "cell-id1"}`,`{"CellId" : "cell-id1","RAT" : "1111","DL_EARFCN" : "1111","PCI" : "111"}`)
    // fmt.Println(err)
    // ls,err := sm.GetCellSet(ctx,"test")
    // fmt.Println(ls)
    // fmt.Println(err)

    err = sm.CreateCellData(ctx, "test",`{"CellId" : "cell-id2"}`,`{"CellId" : "cell-id2","RAT" : "2222","DL_EARFCN" : "2222","PCI" : "222"}`)
    // fmt.Println(err)
    // ls,err = sm.GetCellSet(ctx,"test")
    // fmt.Println(ls)
    // fmt.Println(err)

    err = sm.CreateCellNeighbourRelation(ctx, "test", "nr1", "cell-id1", "cell-id2")
    err = sm.CreateCellNeighbourRelation(ctx, "test", "nr2", "cell-id2", "cell-id2")

    ls,err := sm.GetCellNeighbourRelation(ctx,"test")
    fmt.Println(ls)
    fmt.Println(err)
}


func main() {
	example1()
	example2()
	example3()
}