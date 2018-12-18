package main

import (
	"fmt"
	"log"

	kclient "gopkg.in/jcmturner/gokrb5.v6/client"
	kconfig "gopkg.in/jcmturner/gokrb5.v6/config"
	kcrypto "gopkg.in/jcmturner/gokrb5.v6/crypto"
	"gopkg.in/jcmturner/gokrb5.v6/keytab"
	kmessages "gopkg.in/jcmturner/gokrb5.v6/messages"
	ktypes "gopkg.in/jcmturner/gokrb5.v6/types"
)

func NewKerberosClientWithPassword(username, realm, password string) (*kclient.Client, error) {
	cfg, err := kconfig.Load("/etc/krb5.conf")
	if err != nil {
		return nil, err
	}
	kc := kclient.NewClientWithPassword(username, realm, password)
	kc.WithConfig(cfg)
	return &kc, nil
}

func NewKerberosClientWithKeytab(username, realm string, kt keytab.Keytab) (*kclient.Client, error) {
	cfg, err := kconfig.Load("/etc/krb5.conf")
	if err != nil {
		return nil, fmt.Errorf("error loading kerberos config file: %v", err)
	}
	kc := kclient.NewClientWithKeytab(username, realm, kt)
	kc.WithConfig(cfg)
	return &kc, nil
}

type Client struct {
	kc *kclient.Client
}

func (c *Client) GenerateAuthPayloadFor(servicePrincipal string) (*ClientAuthPayload, error) {
	// https://github.com/jcmturner/gokrb5#generic-kerberos-client
	//
	// To authenticate to a service a client will need to request a service
	// ticket for a Service Principal Name (SPN) and form into an AP_REQ
	// message along with an authenticator encrypted with the session key that
	// was delivered from the KDC along with the service ticket.
	//
	// Get the service ticket and session key for the service the client is
	// authenticating to. The following method will use the client's cache
	// either returning a valid cached ticket, renewing a cached ticket with
	// the KDC or requesting a new ticket from the KDC. Therefore the
	// GetServiceTicket method can be continually used for the most efficient
	// interaction with the KDC.
	tkt, key, err := c.kc.GetServiceTicket(servicePrincipal)
	if err != nil {
		return nil, err
	}
	log.Printf("GetServiceTicket(demo2) got t=[%#v] k=[%#v]", tkt, key)

	// The steps after this will be specific to the application protocol but it
	// will likely involve a client/server Authentication Protocol exchange (AP
	// exchange). This will involve these steps:

	// Generate a new Authenticator and generate a sequence number and subkey:
	auth, err := ktypes.NewAuthenticator(c.kc.Credentials.Realm, c.kc.Credentials.CName)
	if err != nil {
		return nil, err
	}
	etype, err := kcrypto.GetEtype(key.KeyType)
	if err != nil {
		return nil, err
	}
	err = auth.GenerateSeqNumberAndSubKey(key.KeyType, etype.GetKeyByteSize())
	if err != nil {
		return nil, err
	}

	// Set the checksum on the authenticator The checksum is an application
	// specific value. Set as follows:
	// auth.Cksum = ktypes.Checksum{
	// 	CksumType: checksumIDint,
	// 	Checksum:  checksumBytesSlice,
	// }

	// Create the AP_REQ:
	APReq, err := kmessages.NewAPReq(tkt, key, auth)
	if err != nil {
		return nil, err
	}

	return &ClientAuthPayload{APReq: APReq}, nil
}

type ClientAuthPayload struct {
	APReq kmessages.APReq
}
