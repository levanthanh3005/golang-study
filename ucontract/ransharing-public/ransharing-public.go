package ucontract
import (
	"fmt"
	"encoding/json"
	"encoding/base64"

	contractapi "myapp/utility"
	// "github.com/hyperledger/fabric-contract-api-go/contractapi"
	)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

    dataAsBytes := []byte(`{}`)
	err := ctx.GetStub().PutState("MNOs", dataAsBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

    dataAsBytes = []byte(`{}`)
	err = ctx.GetStub().PutState("privateColDisplayName", dataAsBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

    dataAsBytes = []byte(`{"collection": [], "sequence": 1,"version": 1}`)
	err = ctx.GetStub().PutState("privateCollection", dataAsBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

    dataAsBytes = []byte(`{}`)
	err = ctx.GetStub().PutState("InteractionAreaSet", dataAsBytes)

	if err != nil {
		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	return nil
}

func (s *SmartContract) UpdatePrivateCollection(ctx contractapi.TransactionContextInterface, datastr string) (error) {
	fmt.Println("updatePrivateCollection")
	
	var collection map[string]interface{}
	err := StringToMapWithError(datastr, &collection)
	if err != nil {
		return fmt.Errorf("Can not read private collection. %s", err.Error())
	}
	 
	lsCols := (collection["collection"]).([]interface{})
	lsSet,_ := s.ReadMapPublicData(ctx,"InteractionAreaSet")

	fmt.Println(lsSet)
	// for k, _ := range lsSet {
	// 	lsSet[k] = ""
	// }
	fmt.Println("List of collection")
	for _,v := range lsCols {
		vv := v.(map[string]interface{})
		// fmt.Println(vv["interactionArea"])
		str := fmt.Sprintf("%v", vv["interactionArea"])
		// fmt.Println(lsSet[str])
		if lsSet[str] == nil {
			return fmt.Errorf("This interaction does not exist:",str)
		}
	}

	b := []byte(datastr)
	err = ctx.GetStub().PutState("privateCollection", b)

	return nil
}

func (s *SmartContract) CreateNewInteractionArea(ctx contractapi.TransactionContextInterface, areaName string, areaId string) (error) {
	fmt.Println("CreateNewInteractionArea")
	m,err := s.ReadMapPublicData(ctx,"InteractionAreaSet")
	// fmt.Println(m)
	if m[areaId] == nil {
		ssm := make(map[string]string)
		ssm[areaId] = areaName
		m[areaId] = ssm
	}
	// fmt.Println(m)
	content, _ := json.Marshal(m)
	b := []byte(content)
	err = ctx.GetStub().PutState("InteractionAreaSet", b)
	return err
}

func (s *SmartContract) ReadInteractionArea(ctx contractapi.TransactionContextInterface) (string, error) {
	fmt.Println("ReadInteractionArea")
	m,err := s.ReadMapToJSonPublicData(ctx,"InteractionAreaSet")

	return m, err
}

func (s *SmartContract) ReadPrivateCollection(ctx contractapi.TransactionContextInterface) (string, error) {
	fmt.Println("ReadPrivateCollection")
	m,err := s.ReadMapToJSonPublicData(ctx,"privateCollection")

	return m, err
}

func (s *SmartContract) UpdateMNOInfor(ctx contractapi.TransactionContextInterface, mnoData string) (error) {
	// fmt.Println("UpdateMNOInfor")
	mnos,err := s.ReadMapPublicData(ctx,"MNOs")
	// fmt.Println(mnos)
	var mnoObj map[string]interface{}
	err = StringToMapWithError(mnoData, &mnoObj)
	// fmt.Println(mnoObj)

	senderIdSet,_ := s.GetSenderId(ctx)
	// fmt.Println(senderIdSet)
	senderId := senderIdSet["id"]
	// fmt.Println(senderId)

	mnoObj["mspID"] = senderId
	mnos[senderId] = mnoObj

	content, _ := json.Marshal(mnos)
	b := []byte(content)
	err = ctx.GetStub().PutState("MNOs", b)

	return err
	// return nil
}

func (s *SmartContract) GetMNOInfor(ctx contractapi.TransactionContextInterface) (string, error) {
	mnos,err := s.ReadMapPublicData(ctx,"MNOs")
	senderIdSet,_ := s.GetSenderId(ctx)
	content, _ := json.Marshal(mnos[senderIdSet["id"]])
	return string(content),err
}

func (s *SmartContract) GetMNLFullList(ctx contractapi.TransactionContextInterface) (string, error) {
	mnos,err := s.ReadMapPublicData(ctx,"MNOs")
	content, _ := json.Marshal(mnos)
	return string(content),err
}


// // QueryCar returns the car stored in the world state with given id
// func (s *SmartContract) CreatePrivateDataForCar(ctx contractapi.TransactionContextInterface, collection string, carNumber string, privateValue string) (error) {
// 	b := []byte(privateValue)
// 	err := ctx.GetStub().PutPrivateData(collection, carNumber, b)

// 	if err != nil {
// 		return fmt.Errorf("Failed to create private data for car. %s", err.Error())
// 	}

// 	return nil
// }

// func (s *SmartContract) ReadPrivateDataForCar(ctx contractapi.TransactionContextInterface, collection string, carNumber string) (string,error) {
// 	privateValueAsByte, err := ctx.GetStub().GetPrivateData(collection, carNumber)
// 	if err != nil {
// 		return "",fmt.Errorf("Failed to create private data for car. %s", err.Error())
// 	}
// 	return string(privateValueAsByte),nil
// }

func (s *SmartContract) ReadMapPublicData(ctx contractapi.TransactionContextInterface, key string) (map[string]interface{}, error) {
	DataAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	// fmt.Println(string(DataAsBytes))
	b := []byte(DataAsBytes)

	var f interface{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return nil,err
	}
	// fmt.Println(f)
	var m map[string]interface{}
	m = f.(map[string]interface{})

	return m, nil
}

func (s *SmartContract) ReadMapToJSonPublicData(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	DataAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	// fmt.Println(string(DataAsBytes))
	b := []byte(DataAsBytes)

	var f interface{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return "",err
	}
	// fmt.Println(f)
	var m map[string]interface{}
	m = f.(map[string]interface{})
	content, _ := json.Marshal(m)
	return string(content), nil
	// return m, nil
}

func (s *SmartContract) ReadSinglePublicData(ctx contractapi.TransactionContextInterface, key string) (string,error) {
	DataAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	return string(DataAsBytes),nil
}

func StringToMapWithError(objStr string, m *map[string]interface{}) error {
	b := []byte(objStr)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}
	*m = f.(map[string]interface{})
	return nil
}

func (s *SmartContract) GetSenderId(ctx contractapi.TransactionContextInterface) (map[string] string, error) {
	senderId := map[string]string {}
	// ID
	b64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("Failed to read clientID: %v", err)
	}
	// fmt.Println(b64)
	decodeID, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	senderId["id"] = string(decodeID)
	///MSP ID
	b64, err = ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("Failed to read MSP ID: %v", err)
	}

	decodeID, err = base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	senderId["msp"] = string(decodeID)

	///role
	ck := false
	b64, ck, err = ctx.GetClientIdentity().GetAttributeValue("hf.Type")
	if ck == false {
		return nil, fmt.Errorf("Failed to read Role: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to read Role: %v", err)
	}
	decodeID, err = base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	senderId["role"] = string(decodeID)

	return senderId, nil
}


func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}

}