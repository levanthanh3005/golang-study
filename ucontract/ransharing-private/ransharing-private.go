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
	return nil
}

func (s *SmartContract) CreateCellData(ctx contractapi.TransactionContextInterface, privateColName string, publicData string, privateData string) (error) {
	var publicDataMap map[string]interface{}
	err := s.StringToMapWithError(publicData, &publicDataMap)
	if err != nil {
		return fmt.Errorf("Can not parse publicData. %s", err.Error())
	}

	var privateDataMap map[string]interface{}
	err = s.StringToMapWithError(privateData, &privateDataMap)
	if err != nil {
		return fmt.Errorf("Can not parse privateData. %s", err.Error())
	}

	idCell := fmt.Sprintf("%v", (publicDataMap["CellId"]))
	if idCell == "" {
		idCell = fmt.Sprintf("%v", (privateDataMap["CellId"]))
		if idCell == "" {
			return fmt.Errorf("Cell Id is not existed")
		}
	}

	mergeDataMap := make(map[string]interface{})

	for k, v := range publicDataMap {
	    mergeDataMap[k] = v
	}
	for k, v := range privateDataMap {
	    mergeDataMap[k] = v
	}


	cellSet,err := s.GetCellSetMap(ctx, privateColName)
	// fmt.Println(cellSet)
	
	cellList := (cellSet["list"]).(map[string]interface{})
	cellListDetail := (cellSet["listDetail"]).(map[string]interface{})
	cellListStatus := (cellSet["listStatus"]).(map[string]interface{})

	if cellListDetail[idCell] != nil {
		return fmt.Errorf("We already have this cell information")
	}

	cellList[idCell] = idCell
	cellListDetail[idCell] = mergeDataMap

	mergeDataMapList := make(map[string] map[string] interface{})

	mergeDataMapList["publicData"] = publicDataMap
	mergeDataMapList["privateData"] = privateDataMap
	cellListStatus[idCell] = mergeDataMapList

	newCellSet := make(map[string]interface{})
	newCellSet["list"] = cellList
	newCellSet["listDetail"] = cellListDetail
	newCellSet["listStatus"] = cellListStatus

	content, _ := json.Marshal(cellSet)
	b := []byte(string(content))
	err = ctx.GetStub().PutPrivateData(privateColName, "Cells", b)

	return err
}

func (s *SmartContract) UpdateCellData(ctx contractapi.TransactionContextInterface, privateColName string, publicData string, privateData string) (error) {
	var publicDataMap map[string]interface{}
	err := s.StringToMapWithError(publicData, &publicDataMap)
	if err != nil {
		return fmt.Errorf("Can not parse publicData. %s", err.Error())
	}

	var privateDataMap map[string]interface{}
	err = s.StringToMapWithError(privateData, &privateDataMap)
	if err != nil {
		return fmt.Errorf("Can not parse privateData. %s", err.Error())
	}

	idCell := fmt.Sprintf("%v", (publicDataMap["CellId"]))
	if idCell == "" {
		idCell = fmt.Sprintf("%v", (privateDataMap["CellId"]))
		if idCell == "" {
			return fmt.Errorf("Cell Id is not existed")
		}
	}

	mergeDataMap := make(map[string]interface{})

	for k, v := range publicDataMap {
	    mergeDataMap[k] = v
	}
	for k, v := range privateDataMap {
	    mergeDataMap[k] = v
	}


	cellSet,err := s.GetCellSetMap(ctx, privateColName)
	// fmt.Println(cellSet)
	
	cellList := (cellSet["list"]).(map[string]interface{})
	cellListDetail := (cellSet["listDetail"]).(map[string]interface{})
	cellListStatus := (cellSet["listStatus"]).(map[string]interface{})
	if cellListDetail[idCell] == nil {
		return fmt.Errorf("The cell does not existed")
	}

	cellList[idCell] = idCell
	cellListDetail[idCell] = mergeDataMap

	mergeDataMapList := make(map[string] map[string] interface{})

	mergeDataMapList["publicData"] = publicDataMap
	mergeDataMapList["privateData"] = privateDataMap
	cellListStatus[idCell] = mergeDataMapList

	newCellSet := make(map[string]interface{})
	newCellSet["list"] = cellList
	newCellSet["listDetail"] = cellListDetail
	newCellSet["listStatus"] = cellListStatus

	content, _ := json.Marshal(cellSet)
	b := []byte(string(content))
	err = ctx.GetStub().PutPrivateData(privateColName, "Cells", b)

	return err
}

func (s *SmartContract) GetCellSetMap(ctx contractapi.TransactionContextInterface, privateColName string) (map[string]interface{}, error){
	privateValueAsByte, err := s.ReadPrivateDataToBytes(ctx, privateColName,"Cells")
	if len(privateValueAsByte) == 0 {
		privateValueAsByte = []byte(`{"list": {}, "listDetail" : {}, "listStatus": {}}`)
	} 

	cellSet,err := s.DataAsBytesToMap(privateValueAsByte)

	return cellSet,err
}

func (s *SmartContract) GetCellSet(ctx contractapi.TransactionContextInterface, privateColName string) (string, error){
	cellSet, err := s.GetCellSetMap(ctx, privateColName)

	content, _ := json.Marshal(cellSet)
	// fmt.Println("GetCellSet")
	// fmt.Println(content)
	return string(content), err
}

func (s *SmartContract) CreateCellNeighbourRelation(ctx contractapi.TransactionContextInterface, privateColName string, idNR string, idSourceCel string, idTargetCell string) (error) {

	nrSet,err := s.GetCellNeighbourRelationMap(ctx, privateColName)
	if nrSet[idNR] != nil {
		return fmt.Errorf("We already have this NR information")
	}

	cellSet,err := s.GetCellSetMap(ctx, privateColName)
	// fmt.Println(cellSet)
	
	cellListDetail := (cellSet["listDetail"]).(map[string]interface{})
	if cellListDetail[idSourceCel] == nil {
		return fmt.Errorf("The idSourceCel does not existed")
	}
	if cellListDetail[idTargetCell] == nil {
		return fmt.Errorf("The idTargetCell does not existed")
	}

	newNR := make(map[string]string)
	newNR["idNR"] = idNR
	newNR["idSourceCel"] = idSourceCel
	newNR["idTargetCell"] = idTargetCell

	nrList := (nrSet["list"]).(map[string]interface{})
	nrListDetail := (nrSet["listDetail"]).(map[string]interface{})

	nrList[idNR] = idNR
	nrListDetail[idNR] = newNR

	newNRSet := make(map[string]interface{})
	newNRSet["list"] = nrList
	newNRSet["listDetail"] = nrListDetail

	content, _ := json.Marshal(newNRSet)
	b := []byte(string(content))
	err = ctx.GetStub().PutPrivateData(privateColName, "CellNeighbourRelations", b)

	return err
}

func (s *SmartContract) GetCellNeighbourRelation(ctx contractapi.TransactionContextInterface, privateColName string) (string, error) {
	nrSet, err := s.GetCellNeighbourRelationMap(ctx, privateColName)
	content, _ := json.Marshal(nrSet)
	return string(content), err
}

func (s *SmartContract) GetCellNeighbourRelationMap(ctx contractapi.TransactionContextInterface, privateColName string) (map[string]interface{}, error) {
	privateValueAsByte, err := s.ReadPrivateDataToBytes(ctx, privateColName,"CellNeighbourRelations")

	if len(privateValueAsByte) == 0 {
		privateValueAsByte = []byte(`{"list": {}, "listDetail" : {}}`)
	} 

	cellSet,err := s.DataAsBytesToMap(privateValueAsByte)

	return cellSet,err
}

func (s *SmartContract) ReadPrivateDataToBytes(ctx contractapi.TransactionContextInterface, privateColName string, key string) ([]byte, error) {
	privateValueAsByte, err := ctx.GetStub().GetPrivateData(privateColName, key)
	b := []byte(privateValueAsByte)

	return b,err
}

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

func (s *SmartContract) DataAsBytesToString(DataAsBytes []byte) (string,error) {
	return string(DataAsBytes),nil
}

func (s *SmartContract) DataAsBytesToMap(DataAsBytes []byte) (map[string]interface{}, error) {
	// fmt.Println(string(DataAsBytes))
	b := []byte(DataAsBytes)

	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return nil,err
	}
	// fmt.Println(f)
	var m map[string]interface{}
	m = f.(map[string]interface{})

	return m,err
}

func (s *SmartContract) DataAsBytesToJSONString(DataAsBytes []byte) (string, error) {
	m, err := s.DataAsBytesToMap(DataAsBytes)
	content, _ := json.Marshal(m)
	return string(content), err
	// return string(content), err
}

func (s *SmartContract) StringToMapWithError(objStr string, m *map[string]interface{}) error {
	b := []byte(objStr)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}
	*m = f.(map[string]interface{})
	return nil
}

func (s *SmartContract) StringToMapWithError2(objStr string) (map[string]interface{}, error) {
	b := []byte(objStr)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return nil,err
	}
	var m map[string]interface{}
	m = f.(map[string]interface{})
	return m, nil
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