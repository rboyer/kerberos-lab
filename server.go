package main

import (
	"io"
	"log"
	"net/http"
	"os"

	kclient "gopkg.in/jcmturner/gokrb5.v6/client"
	kcredentials "gopkg.in/jcmturner/gokrb5.v6/credentials"
	kservice "gopkg.in/jcmturner/gokrb5.v6/service"
)

type ControlServer struct {
	kc  *kclient.Client
	cfg *kservice.Config
}

func (c *ControlServer) Login(cp *ClientAuthPayload) (*LoginResponse, error) {
	// Now send the AP_REQ to the service. How this is done will be specific to
	// the application use case.
	log.Printf("APReq for demo -> fakeweb is: [%#v]", cp.APReq)
	return nil, io.EOF
}

func (c *ControlServer) Serve() error {
	l := log.New(os.Stderr, "GOKRB5 Service: ", log.Ldate|log.Ltime|log.Lshortfile)

	http.Handle("/login", kservice.SPNEGOKRB5Authenticate(
		http.HandlerFunc(c.HandleLogin),
		c.cfg,
		nil,
	))

	srv := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}

	l.Print("[server] serving /login over 127.0.0.1:8080")
	return srv.ListenAndServe()
}

func (c *ControlServer) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// https://github.com/jcmturner/gokrb5#spnegokerberos-http-service
	log.Printf("[server/http] got login request")

	r.Header.Write(w)
	// const code = http.StatusUnauthorized
	// http.Error(w, http.StatusText(code), code)

	ctx := r.Context()

	creds := ctx.Value(kservice.CTXKeyCredentials).(*kcredentials.Credentials)
	auth := ctx.Value(kservice.CTXKeyAuthenticated).(bool)

	log.Printf("[server/http] got request with creds=[%#v] auth=%v", creds, auth)

}

type LoginResponse struct {
	ServiceAccountUID        string
	ServiceAccountName       string
	ServiceAccountNamespace  string
	ServiceAccountSecretName string
	Role                     string
	Policies                 []string
}
