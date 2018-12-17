package main

import (
	"log"

	kclient "gopkg.in/jcmturner/gokrb5.v6/client"
	kconfig "gopkg.in/jcmturner/gokrb5.v6/config"
	kcrypto "gopkg.in/jcmturner/gokrb5.v6/crypto"
	kmessages "gopkg.in/jcmturner/gokrb5.v6/messages"
	ktypes "gopkg.in/jcmturner/gokrb5.v6/types"
)

func main() {
	cfg, err := kconfig.Load("/etc/krb5.conf")
	if err != nil {
		log.Fatal(err)
	}
	c := kclient.NewClientWithPassword("demo", "KERB.LOCAL", "demo")
	c.WithConfig(cfg)

	err = c.Login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Destroy()

	// ==============
	// 'demo' is going to dial 'demo2'
	// ==============

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
	tkt, key, err := c.GetServiceTicket("demo2")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("GetServiceTicket(demo2) got t=[%#v] k=[%#v]", tkt, key)

	// The steps after this will be specific to the application protocol but it
	// will likely involve a client/server Authentication Protocol exchange (AP
	// exchange). This will involve these steps:

	// Generate a new Authenticator and generate a sequence number and subkey:
	auth, err := ktypes.NewAuthenticator(c.Credentials.Realm, c.Credentials.CName)
	if err != nil {
		log.Fatal(err)
	}
	etype, err := kcrypto.GetEtype(key.KeyType)
	if err != nil {
		log.Fatal(err)
	}
	err = auth.GenerateSeqNumberAndSubKey(key.KeyType, etype.GetKeyByteSize())
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	// Now send the AP_REQ to the service. How this is done will be specific to the application use case.
	log.Printf("APReq for demo -> demo2 is: [%#v]", APReq)
}
