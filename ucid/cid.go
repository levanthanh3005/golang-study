package ucid

import (
	"encoding/base64"
)
type ChaincodeStubInterface interface {
	// GetCreator returns `SignatureHeader.Creator` (e.g. an identity)
	// of the `SignedProposal`. This is the identity of the agent (or user)
	// submitting the transaction.
	GetCreator() ([]byte, error)
}

type ClientIdentity interface {

	// GetID returns the ID associated with the invoking identity.  This ID
	// is guaranteed to be unique within the MSP.
	GetID() (string, error)

	// Return the MSP ID of the client
	GetMSPID() (string, error)

	// GetAttributeValue returns the value of the client's attribute named `attrName`.
	// If the client possesses the attribute, `found` is true and `value` equals the
	// value of the attribute.
	// If the client does not possess the attribute, `found` is false and `value`
	// equals "".
	GetAttributeValue(attrName string) (value string, found bool, err error)

	// AssertAttributeValue verifies that the client has the attribute named `attrName`
	// with a value of `attrValue`; otherwise, an error is returned.
	// AssertAttributeValue(attrName, attrValue string) error

	// GetX509Certificate returns the X509 certificate associated with the client,
	// or nil if it was not identified by an X509 certificate.
	// GetX509Certificate() (*x509.Certificate, error)
}


type ClientID struct {
	stub  ChaincodeStubInterface
	mspID string
	cert  interface{}
	attrs interface{}
}

func GetID(stub ChaincodeStubInterface) (string, error) {
	return "ID",nil
}

// GetMSPID returns the ID of the MSP associated with the identity that
// submitted the transaction
func GetMSPID(stub ChaincodeStubInterface) (string, error) {
	return "MSP",nil
}

// GetAttributeValue returns value of the specified attribute
func GetAttributeValue(stub ChaincodeStubInterface, attrName string) (value string, found bool, err error) {
	return "Attributes",true,nil
}

func (c *ClientID) GetID() (string, error) {
	data := []byte("ClientID")
	str := base64.StdEncoding.EncodeToString(data)
	return str,nil
}

// GetMSPID returns the ID of the MSP associated with the identity that
// submitted the transaction
func (c *ClientID) GetMSPID() (string, error) {
	data := []byte("MSPID")
	str := base64.StdEncoding.EncodeToString(data)
	return str,nil
}

// GetAttributeValue returns value of the specified attribute
func (c *ClientID) GetAttributeValue(attrName string) (value string, found bool, err error) {
	data := []byte("Worker")
	str := base64.StdEncoding.EncodeToString(data)
	return str,true,nil
}

// func (c *ClientID) getIdentity() (*msp.SerializedIdentity, error) {
// 	sid := &msp.SerializedIdentity{}
// 	creator, err := c.stub.GetCreator()
// 	if err != nil || creator == nil {
// 		return nil, fmt.Errorf("failed to get transaction invoker's identity from the chaincode stub: %s", err)
// 	}
// 	err = proto.Unmarshal(creator, sid)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal transaction invoker's identity: %s", err)
// 	}
// 	return sid, nil
// }

// func (c *ClientID) getAttributesFromIdemix() error {
// 	creator, err := c.stub.GetCreator()
// 	attrs, err := attrmgr.New().GetAttributesFromIdemix(creator)
// 	if err != nil {
// 		return fmt.Errorf("failed to get attributes from the transaction invoker's idemix credential: %s", err)
// 	}
// 	c.attrs = attrs
// 	return nil
// }

// New returns an instance of ClientID
func New(stub ChaincodeStubInterface) (*ClientID, error) {
	c := &ClientID{stub: stub}
	err := c.init()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Initialize the client
func (c *ClientID) init() error {
	// signingID, err := c.getIdentity()
	// if err != nil {
	// 	return err
	// }
	// c.mspID = signingID.GetMspid()
	// idbytes := signingID.GetIdBytes()
	// block, _ := pem.Decode(idbytes)
	// if block == nil {
	// 	err := c.getAttributesFromIdemix()
	// 	if err != nil {
	// 		return fmt.Errorf("identity bytes are neither X509 PEM format nor an idemix credential: %s", err)
	// 	}
	// 	return nil
	// }
	// cert, err := x509.ParseCertificate(block.Bytes)
	// if err != nil {
	// 	return fmt.Errorf("failed to parse certificate: %s", err)
	// }
	// c.cert = cert
	// attrs, err := attrmgr.New().GetAttributesFromCert(cert)
	// if err != nil {
	// 	return fmt.Errorf("failed to get attributes from the transaction invoker's certificate: %s", err)
	// }
	// c.attrs = attrs
	c.mspID = "MSPID"
	c.cert = "certificate"
	c.attrs = "attributes"
	return nil
}