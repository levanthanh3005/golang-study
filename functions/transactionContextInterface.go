package functions

// import (
// 	// "strconv"
// 	"myapp/functions"
// 	)


// ClientID holds the information of the transaction creator.
type ClientID struct {
	// ##########
// Take from 
// /Users/vanthanhle/Desktop/Tools/ns3/HyperledgerNetwork/caliper/0403/go/fabric-chaincode-go-master/pkg/cid/cid.go
// ##########
	stub  ChaincodeStubInterface
	mspID string
	cert  *x509.Certificate
	attrs *attrmgr.Attributes
}

// New returns an instance of ClientID
func New(stub ChaincodeStubInterface) (*ClientID, error) {
	c := &ClientID{stub: stub}
	err := c.init()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// GetID returns a unique ID associated with the invoking identity.
func (c *ClientID) GetID() (string, error) {
	// When IdeMix, c.cert is nil for x509 type
	// Here will return "", as there is no x509 type cert for generate id value with logic below.
	if c.cert == nil {
		return "", fmt.Errorf("cannot determine identity")
	}
	// The leading "x509::" distinguishes this as an X509 certificate, and
	// the subject and issuer DNs uniquely identify the X509 certificate.
	// The resulting ID will remain the same if the certificate is renewed.
	id := fmt.Sprintf("x509::%s::%s", getDN(&c.cert.Subject), getDN(&c.cert.Issuer))
	return base64.StdEncoding.EncodeToString([]byte(id)), nil
}

// GetMSPID returns the ID of the MSP associated with the identity that
// submitted the transaction
func (c *ClientID) GetMSPID() (string, error) {
	return c.mspID, nil
}

// GetAttributeValue returns value of the specified attribute
func (c *ClientID) GetAttributeValue(attrName string) (value string, found bool, err error) {
	if c.attrs == nil {
		return "", false, nil
	}
	return c.attrs.Value(attrName)
}

// AssertAttributeValue checks to see if an attribute value equals the specified value
func (c *ClientID) AssertAttributeValue(attrName, attrValue string) error {
	val, ok, err := c.GetAttributeValue(attrName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("attribute '%s' was not found", attrName)
	}
	if val != attrValue {
		return fmt.Errorf("attribute '%s' equals '%s', not '%s'", attrName, val, attrValue)
	}
	return nil
}

// HasOUValue checks if an OU with the specified value is present
func (c *ClientID) HasOUValue(OUValue string) (bool, error) {
	x509Cert := c.cert
	if x509Cert == nil {
		// Here it will return false and an error, as there is no x509 type cert to check for OU values.
		return false, fmt.Errorf("cannot obtain an X509 certificate for the identity")
	}

	for _, OU := range x509Cert.Subject.OrganizationalUnit {
		if OU == OUValue {
			return true, nil
		}
	}
	return false, nil
}

// GetX509Certificate returns the X509 certificate associated with the client,
// or nil if it was not identified by an X509 certificate.
func (c *ClientID) GetX509Certificate() (*x509.Certificate, error) {
	return c.cert, nil
}

// Initialize the client
func (c *ClientID) init() error {
	signingID, err := c.getIdentity()
	if err != nil {
		return err
	}
	c.mspID = signingID.GetMspid()
	idbytes := signingID.GetIdBytes()
	block, _ := pem.Decode(idbytes)
	if block == nil {
		err := c.getAttributesFromIdemix()
		if err != nil {
			return fmt.Errorf("identity bytes are neither X509 PEM format nor an idemix credential: %s", err)
		}
		return nil
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %s", err)
	}
	c.cert = cert
	attrs, err := attrmgr.New().GetAttributesFromCert(cert)
	if err != nil {
		return fmt.Errorf("failed to get attributes from the transaction invoker's certificate: %s", err)
	}
	c.attrs = attrs
	return nil
}

// Unmarshals the bytes returned by ChaincodeStubInterface.GetCreator method and
// returns the resulting msp.SerializedIdentity object
func (c *ClientID) getIdentity() (*msp.SerializedIdentity, error) {
	sid := &msp.SerializedIdentity{}
	creator, err := c.stub.GetCreator()
	if err != nil || creator == nil {
		return nil, fmt.Errorf("failed to get transaction invoker's identity from the chaincode stub: %s", err)
	}
	err = proto.Unmarshal(creator, sid)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction invoker's identity: %s", err)
	}
	return sid, nil
}

func (c *ClientID) getAttributesFromIdemix() error {
	creator, err := c.stub.GetCreator()
	attrs, err := attrmgr.New().GetAttributesFromIdemix(creator)
	if err != nil {
		return fmt.Errorf("failed to get attributes from the transaction invoker's idemix credential: %s", err)
	}
	c.attrs = attrs
	return nil
}

type TransactionContextInterface interface {
	// GetStub should provide a way to access the stub set by Init/Invoke
	GetStub() ChaincodeStubInterface
	// GetClientIdentity should provide a way to access the client identity set by Init/Invoke
	GetClientIdentity() cid.ClientIdentity
}

// SettableTransactionContextInterface defines functions a valid transaction context
// should have. Transaction context's set for contracts to be used in chaincode
// must implement this interface.
type SettableTransactionContextInterface interface {
	// SetStub should provide a way to pass the stub from a chaincode transaction
	// call to the transaction context so that it can be used by contract functions.
	// This is called by Init/Invoke with the stub passed.
	SetStub(shim.ChaincodeStubInterface)
	// SetClientIdentity should provide a way to pass the client identity from a chaincode
	// transaction call to the transaction context so that it can be used by contract functions.
	// This is called by Init/Invoke with the stub passed.
	SetClientIdentity(ci cid.ClientIdentity)
}

// TransactionContext is a basic transaction context to be used in contracts,
// containing minimal required functionality use in contracts as part of
// chaincode. Provides access to the stub and clientIdentity of a transaction.
// If a contract implements the ContractInterface using the Contract struct then
// this is the default transaction context that will be used.
type TransactionContext struct {
	stub           shim.ChaincodeStubInterface
	clientIdentity cid.ClientIdentity
}

// SetStub stores the passed stub in the transaction context
func (ctx *TransactionContext) SetStub(stub shim.ChaincodeStubInterface) {
	ctx.stub = stub
}

// SetClientIdentity stores the passed stub in the transaction context
func (ctx *TransactionContext) SetClientIdentity(ci cid.ClientIdentity) {
	ctx.clientIdentity = ci
}

// GetStub returns the current set stub
func (ctx *TransactionContext) GetStub() shim.ChaincodeStubInterface {
	return ctx.stub
}

// GetClientIdentity returns the current set client identity
func (ctx *TransactionContext) GetClientIdentity() cid.ClientIdentity {
	return ctx.clientIdentity
}



// ##########
// Take from /Users/vanthanhle/Desktop/Tools/ns3/HyperledgerNetwork/caliper/0403/go/fabric-chaincode-go-master/shim/interfaces.go
// ##########

// import (
// 	"strconv"
// 	)


// Chaincode interface must be implemented by all chaincodes. The fabric runs
// the transactions by calling these functions as specified.
type Chaincode interface {
	// Init is called during Instantiate transaction after the chaincode container
	// has been established for the first time, allowing the chaincode to
	// initialize its internal data
	Init(stub ChaincodeStubInterface) string

	// Invoke is called to update or query the ledger in a proposal transaction.
	// Updated state variables are not committed to the ledger until the
	// transaction is committed.
	Invoke(stub ChaincodeStubInterface) string
}

// ChaincodeStubInterface is used by deployable chaincode apps to access and
// modify their ledgers
type ChaincodeStubInterface interface {
	// GetArgs returns the arguments intended for the chaincode Init and Invoke
	// as an array of byte arrays.
	GetArgs() [][]byte

	// GetStringArgs returns the arguments intended for the chaincode Init and
	// Invoke as a string array. Only use GetStringArgs if the client passes
	// arguments intended to be used as strings.
	GetStringArgs() []string

	// GetFunctionAndParameters returns the first argument as the function
	// name and the rest of the arguments as parameters in a string array.
	// Only use GetFunctionAndParameters if the client passes arguments intended
	// to be used as strings.
	GetFunctionAndParameters() (string, []string)

	// GetArgsSlice returns the arguments intended for the chaincode Init and
	// Invoke as a byte array
	GetArgsSlice() ([]byte, error)

	// GetTxID returns the tx_id of the transaction proposal, which is unique per
	// transaction and per client. See
	// https://godoc.org/github.com/hyperledger/fabric-protos-go/common#ChannelHeader
	// for further details.
	GetTxID() string

	// GetChannelID returns the channel the proposal is sent to for chaincode to process.
	// This would be the channel_id of the transaction proposal (see
	// https://godoc.org/github.com/hyperledger/fabric-protos-go/common#ChannelHeader )
	// except where the chaincode is calling another on a different channel.
	GetChannelID() string

	// InvokeChaincode locally calls the specified chaincode `Invoke` using the
	// same transaction context; that is, chaincode calling chaincode doesn't
	// create a new transaction message.
	// If the called chaincode is on the same channel, it simply adds the called
	// chaincode read set and write set to the calling transaction.
	// If the called chaincode is on a different channel,
	// only the Response is returned to the calling chaincode; any PutState calls
	// from the called chaincode will not have any effect on the ledger; that is,
	// the called chaincode on a different channel will not have its read set
	// and write set applied to the transaction. Only the calling chaincode's
	// read set and write set will be applied to the transaction. Effectively
	// the called chaincode on a different channel is a `Query`, which does not
	// participate in state validation checks in subsequent commit phase.
	// If `channel` is empty, the caller's channel is assumed.
	InvokeChaincode(chaincodeName string, args [][]byte, channel string) string

	// GetState returns the value of the specified `key` from the
	// ledger. Note that GetState doesn't read data from the writeset, which
	// has not been committed to the ledger. In other words, GetState doesn't
	// consider data modified by PutState that has not been committed.
	// If the key does not exist in the state database, (nil, nil) is returned.
	GetState(key string) ([]byte, error)

	// PutState puts the specified `key` and `value` into the transaction's
	// writeset as a data-write proposal. PutState doesn't effect the ledger
	// until the transaction is validated and successfully committed.
	// Simple keys must not be an empty string and must not start with a
	// null character (0x00) in order to avoid range query collisions with
	// composite keys, which internally get prefixed with 0x00 as composite
	// key namespace. In addition, if using CouchDB, keys can only contain
	// valid UTF-8 strings and cannot begin with an underscore ("_").
	PutState(key string, value []byte) error

	// DelState records the specified `key` to be deleted in the writeset of
	// the transaction proposal. The `key` and its value will be deleted from
	// the ledger when the transaction is validated and successfully committed.
	DelState(key string) error

	// SetStateValidationParameter sets the key-level endorsement policy for `key`.
	SetStateValidationParameter(key string, ep []byte) error

	// GetStateValidationParameter retrieves the key-level endorsement policy
	// for `key`. Note that this will introduce a read dependency on `key` in
	// the transaction's readset.
	GetStateValidationParameter(key string) ([]byte, error)

	// GetStateByRange returns a range iterator over a set of keys in the
	// ledger. The iterator can be used to iterate over all keys
	// between the startKey (inclusive) and endKey (exclusive).
	// However, if the number of keys between startKey and endKey is greater than the
	// totalQueryLimit (defined in core.yaml), this iterator cannot be used
	// to fetch all keys (results will be capped by the totalQueryLimit).
	// The keys are returned by the iterator in lexical order. Note
	// that startKey and endKey can be empty string, which implies unbounded range
	// query on start or end.
	// Call Close() on the returned StateQueryIteratorInterface object when done.
	// The query is re-executed during validation phase to ensure result set
	// has not changed since transaction endorsement (phantom reads detected).
	// GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)

	// GetStateByRangeWithPagination returns a range iterator over a set of keys in the
	// ledger. The iterator can be used to fetch keys between the startKey (inclusive)
	// and endKey (exclusive).
	// When an empty string is passed as a value to the bookmark argument, the returned
	// iterator can be used to fetch the first `pageSize` keys between the startKey
	// (inclusive) and endKey (exclusive).
	// When the bookmark is a non-emptry string, the iterator can be used to fetch
	// the first `pageSize` keys between the bookmark (inclusive) and endKey (exclusive).
	// Note that only the bookmark present in a prior page of query results (ResponseMetadata)
	// can be used as a value to the bookmark argument. Otherwise, an empty string must
	// be passed as bookmark.
	// The keys are returned by the iterator in lexical order. Note
	// that startKey and endKey can be empty string, which implies unbounded range
	// query on start or end.
	// Call Close() on the returned StateQueryIteratorInterface object when done.
	// This call is only supported in a read only transaction.
	// GetStateByRangeWithPagination(startKey, endKey string, pageSize int32,
	// 	bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)

	// GetStateByPartialCompositeKey queries the state in the ledger based on
	// a given partial composite key. This function returns an iterator
	// which can be used to iterate over all composite keys whose prefix matches
	// the given partial composite key. However, if the number of matching composite
	// keys is greater than the totalQueryLimit (defined in core.yaml), this iterator
	// cannot be used to fetch all matching keys (results will be limited by the totalQueryLimit).
	// The `objectType` and attributes are expected to have only valid utf8 strings and
	// should not contain U+0000 (nil byte) and U+10FFFF (biggest and unallocated code point).
	// See related functions SplitCompositeKey and CreateCompositeKey.
	// Call Close() on the returned StateQueryIteratorInterface object when done.
	// The query is re-executed during validation phase to ensure result set
	// has not changed since transaction endorsement (phantom reads detected).
	// GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error)

	// GetStateByPartialCompositeKeyWithPagination queries the state in the ledger based on
	// a given partial composite key. This function returns an iterator
	// which can be used to iterate over the composite keys whose
	// prefix matches the given partial composite key.
	// When an empty string is passed as a value to the bookmark argument, the returned
	// iterator can be used to fetch the first `pageSize` composite keys whose prefix
	// matches the given partial composite key.
	// When the bookmark is a non-emptry string, the iterator can be used to fetch
	// the first `pageSize` keys between the bookmark (inclusive) and the last matching
	// composite key.
	// Note that only the bookmark present in a prior page of query result (ResponseMetadata)
	// can be used as a value to the bookmark argument. Otherwise, an empty string must
	// be passed as bookmark.
	// The `objectType` and attributes are expected to have only valid utf8 strings
	// and should not contain U+0000 (nil byte) and U+10FFFF (biggest and unallocated
	// code point). See related functions SplitCompositeKey and CreateCompositeKey.
	// Call Close() on the returned StateQueryIteratorInterface object when done.
	// This call is only supported in a read only transaction.
	// GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string,
	// 	pageSize int32, bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)

	// CreateCompositeKey combines the given `attributes` to form a composite
	// key. The objectType and attributes are expected to have only valid utf8
	// strings and should not contain U+0000 (nil byte) and U+10FFFF
	// (biggest and unallocated code point).
	// The resulting composite key can be used as the key in PutState().
	CreateCompositeKey(objectType string, attributes []string) (string, error)

	// SplitCompositeKey splits the specified key into attributes on which the
	// composite key was formed. Composite keys found during range queries
	// or partial composite key queries can therefore be split into their
	// composite parts.
	SplitCompositeKey(compositeKey string) (string, []string, error)

	// GetQueryResult performs a "rich" query against a state database. It is
	// only supported for state databases that support rich query,
	// e.g.CouchDB. The query string is in the native syntax
	// of the underlying state database. An iterator is returned
	// which can be used to iterate over all keys in the query result set.
	// However, if the number of keys in the query result set is greater than the
	// totalQueryLimit (defined in core.yaml), this iterator cannot be used
	// to fetch all keys in the query result set (results will be limited by
	// the totalQueryLimit).
	// The query is NOT re-executed during validation phase, phantom reads are
	// not detected. That is, other committed transactions may have added,
	// updated, or removed keys that impact the result set, and this would not
	// be detected at validation/commit time.  Applications susceptible to this
	// should therefore not use GetQueryResult as part of transactions that update
	// ledger, and should limit use to read-only chaincode operations.
	GetQueryResult(query string) (string, error)

	// GetQueryResultWithPagination performs a "rich" query against a state database.
	// It is only supported for state databases that support rich query,
	// e.g., CouchDB. The query string is in the native syntax
	// of the underlying state database. An iterator is returned
	// which can be used to iterate over keys in the query result set.
	// When an empty string is passed as a value to the bookmark argument, the returned
	// iterator can be used to fetch the first `pageSize` of query results.
	// When the bookmark is a non-emptry string, the iterator can be used to fetch
	// the first `pageSize` keys between the bookmark and the last key in the query result.
	// Note that only the bookmark present in a prior page of query results (ResponseMetadata)
	// can be used as a value to the bookmark argument. Otherwise, an empty string
	// must be passed as bookmark.
	// This call is only supported in a read only transaction.
	// GetQueryResultWithPagination(query string, pageSize int32,
	// 	bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)

	// GetHistoryForKey returns a history of key values across time.
	// For each historic key update, the historic value and associated
	// transaction id and timestamp are returned. The timestamp is the
	// timestamp provided by the client in the proposal header.
	// GetHistoryForKey requires peer configuration
	// core.ledger.history.enableHistoryDatabase to be true.
	// The query is NOT re-executed during validation phase, phantom reads are
	// not detected. That is, other committed transactions may have updated
	// the key concurrently, impacting the result set, and this would not be
	// detected at validation/commit time. Applications susceptible to this
	// should therefore not use GetHistoryForKey as part of transactions that
	// update ledger, and should limit use to read-only chaincode operations.
	// GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)

	// GetPrivateData returns the value of the specified `key` from the specified
	// `collection`. Note that GetPrivateData doesn't read data from the
	// private writeset, which has not been committed to the `collection`. In
	// other words, GetPrivateData doesn't consider data modified by PutPrivateData
	// that has not been committed.
	GetPrivateData(collection, key string) ([]byte, error)

	// GetPrivateDataHash returns the hash of the value of the specified `key` from the specified
	// `collection`
	GetPrivateDataHash(collection, key string) ([]byte, error)

	// PutPrivateData puts the specified `key` and `value` into the transaction's
	// private writeset. Note that only hash of the private writeset goes into the
	// transaction proposal response (which is sent to the client who issued the
	// transaction) and the actual private writeset gets temporarily stored in a
	// transient store. PutPrivateData doesn't effect the `collection` until the
	// transaction is validated and successfully committed. Simple keys must not
	// be an empty string and must not start with a null character (0x00) in order
	// to avoid range query collisions with composite keys, which internally get
	// prefixed with 0x00 as composite key namespace. In addition, if using
	// CouchDB, keys can only contain valid UTF-8 strings and cannot begin with an
	// an underscore ("_").
	PutPrivateData(collection string, key string, value []byte) error

	// DelPrivateData records the specified `key` to be deleted in the private writeset
	// of the transaction. Note that only hash of the private writeset goes into the
	// transaction proposal response (which is sent to the client who issued the
	// transaction) and the actual private writeset gets temporarily stored in a
	// transient store. The `key` and its value will be deleted from the collection
	// when the transaction is validated and successfully committed.
	DelPrivateData(collection, key string) error

	// SetPrivateDataValidationParameter sets the key-level endorsement policy
	// for the private data specified by `key`.
	SetPrivateDataValidationParameter(collection, key string, ep []byte) error

	// GetPrivateDataValidationParameter retrieves the key-level endorsement
	// policy for the private data specified by `key`. Note that this introduces
	// a read dependency on `key` in the transaction's readset.
	GetPrivateDataValidationParameter(collection, key string) ([]byte, error)

	// GetPrivateDataByRange returns a range iterator over a set of keys in a
	// given private collection. The iterator can be used to iterate over all keys
	// between the startKey (inclusive) and endKey (exclusive).
	// The keys are returned by the iterator in lexical order. Note
	// that startKey and endKey can be empty string, which implies unbounded range
	// query on start or end.
	// Call Close() on the returned StateQueryIteratorInterface object when done.
	// The query is re-executed during validation phase to ensure result set
	// has not changed since transaction endorsement (phantom reads detected).
	// GetPrivateDataByRange(collection, startKey, endKey string) (StateQueryIteratorInterface, error)

	// GetPrivateDataByPartialCompositeKey queries the state in a given private
	// collection based on a given partial composite key. This function returns
	// an iterator which can be used to iterate over all composite keys whose prefix
	// matches the given partial composite key. The `objectType` and attributes are
	// expected to have only valid utf8 strings and should not contain
	// U+0000 (nil byte) and U+10FFFF (biggest and unallocated code point).
	// See related functions SplitCompositeKey and CreateCompositeKey.
	// Call Close() on the returned StateQueryIteratorInterface object when done.
	// The query is re-executed during validation phase to ensure result set
	// has not changed since transaction endorsement (phantom reads detected).
	// GetPrivateDataByPartialCompositeKey(collection, objectType string, keys []string) (StateQueryIteratorInterface, error)

	// GetPrivateDataQueryResult performs a "rich" query against a given private
	// collection. It is only supported for state databases that support rich query,
	// e.g.CouchDB. The query string is in the native syntax
	// of the underlying state database. An iterator is returned
	// which can be used to iterate (next) over the query result set.
	// The query is NOT re-executed during validation phase, phantom reads are
	// not detected. That is, other committed transactions may have added,
	// updated, or removed keys that impact the result set, and this would not
	// be detected at validation/commit time.  Applications susceptible to this
	// should therefore not use GetPrivateDataQueryResult as part of transactions that update
	// ledger, and should limit use to read-only chaincode operations.
	// GetPrivateDataQueryResult(collection, query string) (StateQueryIteratorInterface, error)

	// GetCreator returns `SignatureHeader.Creator` (e.g. an identity)
	// of the `SignedProposal`. This is the identity of the agent (or user)
	// submitting the transaction.
	GetCreator() ([]byte, error)

	// GetTransient returns the `ChaincodeProposalPayload.Transient` field.
	// It is a map that contains data (e.g. cryptographic material)
	// that might be used to implement some form of application-level
	// confidentiality. The contents of this field, as prescribed by
	// `ChaincodeProposalPayload`, are supposed to always
	// be omitted from the transaction and excluded from the ledger.
	GetTransient() (map[string][]byte, error)

	// GetBinding returns the transaction binding, which is used to enforce a
	// link between application data (like those stored in the transient field
	// above) to the proposal itself. This is useful to avoid possible replay
	// attacks.
	GetBinding() ([]byte, error)

	// GetDecorations returns additional data (if applicable) about the proposal
	// that originated from the peer. This data is set by the decorators of the
	// peer, which append or mutate the chaincode input passed to the chaincode.
	GetDecorations() map[string][]byte

	// GetSignedProposal returns the SignedProposal object, which contains all
	// data elements part of a transaction proposal.
	// GetSignedProposal() (*pb.SignedProposal, error)

	// GetTxTimestamp returns the timestamp when the transaction was created. This
	// is taken from the transaction ChannelHeader, therefore it will indicate the
	// client's timestamp and will have the same value across all endorsers.
	// GetTxTimestamp() (*timestamp.Timestamp, error)

	// SetEvent allows the chaincode to set an event on the response to the
	// proposal to be included as part of a transaction. The event will be
	// available within the transaction in the committed block regardless of the
	// validity of the transaction.
	// Only a single event can be included in a transaction, and must originate
	// from the outer-most invoked chaincode in chaincode-to-chaincode scenarios.
	// The marshaled ChaincodeEvent will be available in the transaction's ChaincodeAction.events field.
	SetEvent(name string, payload []byte) error
}