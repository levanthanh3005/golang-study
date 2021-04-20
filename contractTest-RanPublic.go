package main
import (
	"fmt"
	// "encoding/json"
	// "strconv"
	contract "myapp/ucontract/ransharing-public"
	contractapi "myapp/utility"
	)

func example0() {
	// fmt.Println(contractapi.stringToJsonString("{test}"))
	m := contractapi.StringToMap(`{"test":"1"}`)
	s,_ := contractapi.MapToJsonString(m)
	fmt.Println(s)
}
func example1() {
    sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)
    stub.InitState()
    ctx.SetStub(stub)
    
    sm.InitLedger(ctx)

    fmt.Println(sm.GetSenderId(ctx))
    // fmt.Println()
}

func example2() {
	fmt.Println("Update and Read PrivateCollection & InteractionArea")
	sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)
    stub.InitState()
    ctx.SetStub(stub)
    
    sm.InitLedger(ctx)
    fmt.Println(sm.GetSenderId(ctx))
    sm.CreateNewInteractionArea(ctx, "Brenero", "Brenero")
    sm.CreateNewInteractionArea(ctx, "Italo-Swiss", "Ita-Sw")

	InteractionAreaSet,_ := sm.ReadInteractionArea(ctx)
	// fmt.Println("InteractionAreaSet")
	fmt.Println(InteractionAreaSet)

	var obj = 
	`
	{
	"collection": [
	  {
	     "name": "Manager-0",
	     "policy": "OR('Org1MSP.member', 'Org2MSP.member')",
	     "requiredPeerCount": 0,
	     "maxPeerCount": 3,
	     "blockToLive": 0,
	     "interactionArea" : "Brenero",
	     "memberOnlyRead": true
	  },
	  {
	     "name": "TIMOnly-0",
	     "policy": "OR('Org10MSP.member')",
	     "requiredPeerCount": 0,
	     "maxPeerCount": 3,
	     "blockToLive": 0,
	         "interactionArea" : "Ita-Sw",
	     "memberOnlyRead": true
	  }
	]
	}
	`
	err := sm.UpdatePrivateCollection(ctx, obj)
	fmt.Println(err)
	privateCollection,_ := sm.ReadPrivateCollection(ctx)
	fmt.Println(privateCollection);
	// await myContract.updatePrivateCollection(ctx, JSON.stringify(obj));
}

func example3() {
	fmt.Println("Update and Read MNO Infor")
	sm := new(contract.SmartContract)
    ctx := new(contractapi.TransactionContext)
    stub := new(contractapi.ChaincodeStub)
    stub.InitState()
    ctx.SetStub(stub)
    
    sm.InitLedger(ctx)
    fmt.Println(sm.GetSenderId(ctx))
    sm.UpdateMNOInfor(ctx, `{"MNOName" : "DeutschTelekom", "PLMNId" : "56259", "MNOHost": "10.10.21.109"}`)
    mnoInfor,_:=sm.GetMNOInfor(ctx)
    fmt.Println(mnoInfor)
    mnos,_:=sm.GetMNLFullList(ctx)
    fmt.Println(mnos)
}

func main() {
	example2()
	example3()
}