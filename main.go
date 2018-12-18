package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"gopkg.in/jcmturner/gokrb5.v6/keytab"
	kservice "gopkg.in/jcmturner/gokrb5.v6/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var srv *ControlServer
	{
		kt, err := keytab.Load("/keytabs/fakeweb.keytab")
		if err != nil {
			return fmt.Errorf("[server] kclient parse keytab error: %v", err)
		}

		// this client for the service may not be necessary
		kc, err := NewKerberosClientWithKeytab("HTTP/fakeweb", "KERB.LOCAL", kt)
		if err != nil {
			return fmt.Errorf("[server] kclient create error: %v", err)
		}

		err = kc.Login()
		if err != nil {
			return fmt.Errorf("[server] kclient login error: %v", err)
		}
		defer kc.Destroy()

		srv = &ControlServer{
			kc:  kc,
			cfg: kservice.NewConfig(kt),
		}

		go func() {
			if err := srv.Serve(); err != nil {
				log.Fatal(err)
			}
		}()
	}
	_ = srv

	var cli *Client
	{
		kc, err := NewKerberosClientWithPassword("demo", "KERB.LOCAL", "demo")
		if err != nil {
			return fmt.Errorf("[client] kclient create error: %v", err)
		}

		err = kc.Login()
		if err != nil {
			return fmt.Errorf("[client] kclient login error: %v", err)
		}
		defer kc.Destroy()

		cli = &Client{kc: kc}
	}

	// cp, err := cli.GenerateAuthPayloadFor("fakeweb")
	// if err != nil {
	// 	return fmt.Errorf("[client] kclient payload generation error: %v", err)
	// }

	// resp, err := srv.Login(cp)
	// if err != nil {
	// 	return fmt.Errorf("[client] login failed: %v", err)
	// }
	// _ = resp

	{
		data := strings.NewReader(`{"placeholder":true}`)
		r, err := http.NewRequest("POST", "http://127.0.0.1:8080/login", data)
		if err != nil {
			return fmt.Errorf("[client] http request construction error: %v", err)
		}

		r.Header.Set("content-type", "application/json")

		err = cli.kc.SetSPNEGOHeader(r, "HTTP/fakeweb")
		if err != nil {
			return fmt.Errorf("[client] http request SPNEGO header set error: %v", err)
		}

		hc := cleanhttp.DefaultClient()

		res, err := hc.Do(r)
		if err != nil {
			return fmt.Errorf("[client] http Do error: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return fmt.Errorf("bogus status: got %v", res.Status)
		}

		got, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("[client] http body drain error: %v", err)
		}
		log.Printf("[client] got data: [%s]", string(got))

	}

	return nil
}
