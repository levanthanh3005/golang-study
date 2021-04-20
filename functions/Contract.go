package functions

import (
	"strconv"
	)

func GetValue(a int) string{
    return "Hello from this another package "+strconv.Itoa(a)
}

type Contract struct {
	Name                      string
	Info                      string
	UnknownTransaction        interface{}
	BeforeTransaction         interface{}
	AfterTransaction          interface{}
	TransactionContextHandler string
}

// GetInfo returns the info about the contract for use in metadata
func (c *Contract) GetInfo() string {
	return c.Info
}

// GetUnknownTransaction returns the current set unknownTransaction, may be nil
func (c *Contract) GetUnknownTransaction() interface{} {
	return c.UnknownTransaction
}

// GetBeforeTransaction returns the current set beforeTransaction, may be nil
func (c *Contract) GetBeforeTransaction() interface{} {
	return c.BeforeTransaction
}

// GetAfterTransaction returns the current set afterTransaction, may be nil
func (c *Contract) GetAfterTransaction() interface{} {
	return c.AfterTransaction
}

// GetName returns the name of the contract
func (c *Contract) GetName() string {
	return c.Name
}

type ContractInterface interface {
	GetInfo() string
	GetUnknownTransaction() interface{}
	GetBeforeTransaction() interface{}
	GetAfterTransaction() interface{}
	GetName() string
	GetTransactionContextHandler() string
}


func (c *Contract) GetTransactionContextHandler() string {

	return "c.TransactionContextHandler"
}

type contractChaincodeContract struct {
	info                      string
}

type ContractChaincode struct {
	DefaultContract       string
	contracts             map[string]contractChaincodeContract
	metadata              string
	Info                  string
	TransactionSerializer string
}

func NewChaincode(contracts ...ContractInterface) (*ContractChaincode, error) {
	// ciMethods := getCiMethods()

	// cc := new(ContractChaincode)
	// cc.contracts = make(map[string]contractChaincodeContract)

	// for _, contract := range contracts {
	// 	additionalExcludes := []string{}
	// 	if castContract, ok := contract.(IgnoreContractInterface); ok {
	// 		additionalExcludes = castContract.GetIgnoredFunctions()
	// 	}

	// 	err := cc.addContract(contract, append(ciMethods, additionalExcludes...))

	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// sysC := new(SystemContract)
	// sysC.Name = SystemContractName

	// cc.addContract(sysC, ciMethods) // should never error as system contract is good

	// err := cc.augmentMetadata()

	// if err != nil {
	// 	return nil, err
	// }

	// metadataJSON, _ := json.Marshal(cc.metadata)

	// sysC.setMetadata(string(metadataJSON))

	// cc.TransactionSerializer = new(serializer.JSONSerializer)

	return nil, nil
}

// Start starts the chaincode in the fabric shim
func (cc *ContractChaincode) Start() error {
	return nil
}