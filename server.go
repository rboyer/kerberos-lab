package main

import (
	"io"
	"log"
	"net/http"

	kclient "gopkg.in/jcmturner/gokrb5.v6/client"
)

type ControlServer struct {
	kc *kclient.Client
}

func (c *ControlServer) Login(cp *ClientAuthPayload) (*LoginResponse, error) {
	// Now send the AP_REQ to the service. How this is done will be specific to
	// the application use case.
	log.Printf("APReq for demo -> fakeweb is: [%#v]", cp.APReq)
	return nil, io.EOF
}

func (c *ControlServer) Serve() error {
	http.HandleFunc("/login", c.HandleLogin)
	srv := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}
	log.Print("[server] serving /login over 127.0.0.1:8080")
	return srv.ListenAndServe()
}

func (c *ControlServer) HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("[server/http] got login request")
	r.Header.Write(w)
	// const code = http.StatusUnauthorized
	// http.Error(w, http.StatusText(code), code)
}

type LoginResponse struct {
	ServiceAccountUID        string
	ServiceAccountName       string
	ServiceAccountNamespace  string
	ServiceAccountSecretName string
	Role                     string
	Policies                 []string
}
