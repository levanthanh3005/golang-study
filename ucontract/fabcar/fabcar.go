package ucontract
import (
	"fmt"
	"encoding/json"
	"strconv"
	"encoding/base64"

	contractapi "myapp/utility"
	// "github.com/hyperledger/fabric-contract-api-go/contractapi"
	)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a car
type Car struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Car
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Car{
		Car{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		Car{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
		Car{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		Car{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
		Car{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
		Car{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		Car{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Car{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Car{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		Car{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}

	for i, car := range cars {
		carAsBytes, _ := json.Marshal(car)
		err := ctx.GetStub().PutState("CAR"+strconv.Itoa(i), carAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}


// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, carNumber string, make string, model string, colour string, owner string) error {
	car := Car{
		Make:   make,
		Model:  model,
		Colour: colour,
		Owner:  owner,
	}

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, carNumber string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carNumber)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) RemoveCar(ctx contractapi.TransactionContextInterface, carNumber string) (error) {
	err := ctx.GetStub().DelState(carNumber)

	if err != nil {
		return fmt.Errorf("Failed to delete from world state. %s", err.Error())
	}

	return nil
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) CreatePrivateDataForCar(ctx contractapi.TransactionContextInterface, collection string, carNumber string, privateValue string) (error) {
	b := []byte(privateValue)
	err := ctx.GetStub().PutPrivateData(collection, carNumber, b)

	if err != nil {
		return fmt.Errorf("Failed to create private data for car. %s", err.Error())
	}

	return nil
}

func (s *SmartContract) ReadPrivateDataForCar(ctx contractapi.TransactionContextInterface, collection string, carNumber string) (string,error) {
	privateValueAsByte, err := ctx.GetStub().GetPrivateData(collection, carNumber)
	if err != nil {
		return "",fmt.Errorf("Failed to create private data for car. %s", err.Error())
	}
	return string(privateValueAsByte),nil
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

// async getSenderId(ctx){
//         // https://hyperledger-fabric.readthedocs.io/en/release-2.0/developapps/transactioncontext.html#structure
//         //https://stackoverflow.com/questions/57772249/get-mspid-from-a-command-line-request-in-chaincode
//         let cid = new ClientIdentity(ctx.stub);
//         let mspID = cid.getMSPID();//Org1MSP
//         let id = cid.getID();
// //   "id": "x509::/OU=client/OU=org1/OU=department1/CN=appUser::/C=US/ST=North Carolina/L=Durham/O=org1.example.com/CN=ca.org1.example.com"

//         let affiliation = cid.getAttributeValue("hf.Affiliation");
//         let enrollmentID = cid.getAttributeValue("hf.EnrollmentID");
//         let role = cid.getAttributeValue("hf.Type");

//         var beginId = id.indexOf("/CN=");
//         var endId = id.lastIndexOf("::/C=");
//         let userId = id.substring(beginId + 4, endId);
//         return {msp : mspID, id: id, userId:userId, affiliation : affiliation, enrollmentID : enrollmentID, role : role};
//     }


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