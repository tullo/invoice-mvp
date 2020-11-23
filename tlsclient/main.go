package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Read the CA signed certificate.
	cert, err := ioutil.ReadFile("../localhost+2.pem")
	if err != nil {
		log.Fatal("Couldn't load certificate file", err)
	}
	// Add cert to the pool of trusted certificates.
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	// Use trusted server certificate for client transport config.
	var c tls.Config
	c.RootCAs = certPool
	var t http.Transport
	t.TLSClientConfig = &c

	var client http.Client
	client.Transport = &t

	// TLS-based authentification. Client makes sure to only talk to servers
	// in possession of the private key used to sign the certificate.
	req, _ := http.NewRequest("GET", "https://127.0.0.1:8443/activities", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8")
	res, err := client.Do(req)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Status", res.Status)
}
