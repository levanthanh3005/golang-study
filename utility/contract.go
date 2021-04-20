// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package utility

import (
	"fmt"
	// "encoding/json"
	cid "myapp/ucid"
	)

type Contract struct {
	Name                      string
	Info                      string
	UnknownTransaction        interface{}
	BeforeTransaction         interface{}
	AfterTransaction          interface{}
	TransactionContextHandler string
}

type ContractChaincode struct {
	DefaultContract       string
	contracts             map[string]interface{}
	metadata              interface{}
	Info                  interface{}
	TransactionSerializer interface{}
}

// type ClientIdentity interface{}

type ContractInterface interface {
	GetInfo() interface{}
	GetUnknownTransaction() interface{}

	GetBeforeTransaction() interface{}

	GetAfterTransaction() interface{}
	GetName() string
	GetTransactionContextHandler() interface{}
}

type TransactionContextInterface interface {
	// GetStub should provide a way to access the stub set by Init/Invoke
	GetStub() ChaincodeStubInterface
	// GetClientIdentity should provide a way to access the client identity set by Init/Invoke
	GetClientIdentity() cid.ClientIdentity
}

type TransactionContext struct {
	stub           ChaincodeStubInterface
	clientIdentity cid.ClientIdentity
}


// SetStub stores the passed stub in the transaction context
func (ctx *TransactionContext) SetStub(stub ChaincodeStubInterface) {
	ctx.stub = stub
}

// SetClientIdentity stores the passed stub in the transaction context
func (ctx *TransactionContext) SetClientIdentity(ci cid.ClientIdentity) {
	ctx.clientIdentity = ci
}

// GetStub returns the current set stub
func (ctx *TransactionContext) GetStub() ChaincodeStubInterface {
	return ctx.stub
}

// GetClientIdentity returns the current set client identity
func (ctx *TransactionContext) GetClientIdentity() cid.ClientIdentity {
	c := new(cid.ClientID)
	ctx.clientIdentity = c
	return ctx.clientIdentity
}

type ChaincodeStubInterface interface {
	// GetArgs() [][]byte
	// GetStringArgs() []string

	// GetFunctionAndParameters() (string, []string)

	// GetArgsSlice() ([]byte, error)
	GetTxID() string

	// GetChannelID() string

	// InvokeChaincode(chaincodeName string, args [][]byte, channel string) string

	GetState(key string) ([]byte, error)

	PutState(key string, value []byte) error

	DelState(key string) error

	// SetStateValidationParameter(key string, ep []byte) error

	// GetStateValidationParameter(key string) ([]byte, error)

	// CreateCompositeKey(objectType string, attributes []string) (string, error)

	// SplitCompositeKey(compositeKey string) (string, []string, error)

	// GetQueryResult(query string) (string, error)
	GetPrivateData(collection, key string) ([]byte, error)

	// GetPrivateDataHash(collection, key string) ([]byte, error)

	PutPrivateData(collection string, key string, value []byte) error

	DelPrivateData(collection, key string) error

	// SetPrivateDataValidationParameter(collection, key string, ep []byte) error
	// GetPrivateDataValidationParameter(collection, key string) ([]byte, error)
	// GetCreator() ([]byte, error)
	// GetTransient() (map[string][]byte, error)
	// GetBinding() ([]byte, error)
	// GetDecorations() map[string][]byte
	// SetEvent(name string, payload []byte) error

	GetAllState() map[string][]byte
}

type ChaincodeStub struct {
	TxID                       string
	ChannelID                  string
	// chaincodeEvent             *pb.ChaincodeEvent
	args                       [][]byte
	// handler                    *Handler
	// signedProposal             *pb.SignedProposal
	// proposal                   *pb.Proposal
	validationParameterMetakey string

	// Additional fields extracted from the signedProposal
	creator   []byte
	transient map[string][]byte
	binding   []byte

	decorations map[string][]byte
	publicData map[string][]byte
	privateData map[string]map[string][]byte

}

func (s *ChaincodeStub) InitState() (error) {
	s.publicData = map[string] []byte {}
	s.privateData = map[string]map[string][]byte {}
	return nil
}

func (s *ChaincodeStub) GetTxID() (string) {
	// Access public data by setting the collection to empty string
	// collection := ""
	return "TX1234"
}

func (s *ChaincodeStub) DelState(key string) error {
	delete(s.publicData,key)
	return nil
}

func (s *ChaincodeStub) GetState(key string) ([]byte, error) {

	return s.publicData[key], nil
}

func (s *ChaincodeStub) GetAllState() (map[string][]byte) {

	return s.publicData
}

func (s *ChaincodeStub) PutState(key string, value []byte) error {
	s.publicData[key] = value

	return (nil)
}

func (s *ChaincodeStub) PutPrivateData(collection string, key string, value []byte) error {
	if s.privateData[collection] == nil {
		s.privateData[collection] = map[string] []byte {}
	}

	s.privateData[collection][key] = value

	return nil
}

func (s *ChaincodeStub) DelPrivateData(collection, key string) error {
	if s.privateData[collection] == nil {
		return fmt.Errorf("The private collection does not existed")
	}
	delete(s.privateData[collection],key)
	return nil
}

func (s *ChaincodeStub) GetPrivateData(collection, key string) ([]byte, error) {
	if s.privateData[collection] == nil {
		return nil, fmt.Errorf("The private collection does not existed")
	}
	return s.privateData[collection][key],nil
}

// GetInfo returns the info about the contract for use in metadata
func (c *Contract) GetInfo() interface{} {
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

// GetTransactionContextHandler returns the current transaction context set for
// the contract. If none has been set then TransactionContext will be returned
func (c *Contract) GetTransactionContextHandler() interface{} {
	return c.TransactionContextHandler
}


func NewChaincode(contracts ...ContractInterface) (*ContractChaincode, error) {
	// ciMethods := getCiMethods()

	cc := new(ContractChaincode)
	cc.contracts = make(map[string]interface{})

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

func (c *ContractChaincode) Start() (error) {
	return nil
}