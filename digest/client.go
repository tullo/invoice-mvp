package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	body := `{"name": "3skills"}`
	digestPost("https://127.0.0.1:8443", "/customers", []byte(body))
}

func digestParts(resp *http.Response) map[string]string {
	m := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		parts := []string{"algorithm", "nonce", "realm", "qop"}
		auth := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range auth {
			for _, part := range parts {
				if strings.Contains(r, part) {
					m[part] = strings.Split(r, `"`)[1]
				}
			}
		}
	}
	return m
}

func digestPost(host string, uri string, payload []byte) {
	url := host + uri
	method := "POST"

	// send initial request
	r, err := http.NewRequest(method, url, nil)
	r.Header.Set("Content-Type", "application/json")
	var client http.Client
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		log.Printf("Unexpected status code %q", resp.StatusCode)
		return

	}
	// server auth challenge response
	m := digestParts(resp)
	m["uri"] = uri
	m["method"] = method
	m["username"] = "go"
	m["password"] = "time"
	r, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	r.Header.Set("Authorization", digestAuth(m))
	r.Header.Set("Content-Type", "application/json")

	// send post request with valid auth values
	resp, err = client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Println("Customer NOT created!", resp.StatusCode)
		return
	}

	log.Println("Customer created!")
}

func digestAuth(m map[string]string) string {
	var alg = "MD5"
	if a, ok := m["algorithm"]; ok {
		switch a {
		case "SHA-256":
			alg = "SHA-256"
		}
	}
	h1 := hashsum(fmt.Sprintf("%s:%s:%s", m["username"], m["realm"], m["password"]), alg)
	h2 := hashsum(fmt.Sprintf("%s:%s", m["method"], m["uri"]), alg)
	nonceCount := 4567
	cnonce := cnonce()
	response := hashsum(fmt.Sprintf("%s:%s:%08x:%s:%s:%s", h1, m["nonce"], nonceCount, cnonce, m["qop"], h2), alg)
	format := `Digest username=%q, realm="%s", nonce=%q, uri=%q, cnonce=%q, nc="%08x", qop=%q, response=%q`
	return fmt.Sprintf(format, m["username"], m["realm"], m["nonce"], m["uri"], cnonce, nonceCount, m["qop"], response)
}

func cnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}

func hashsum(s, alg string) string {
	var h hash.Hash
	switch alg {
	case "MD5":
		h = md5.New()
	case "SHA-256":
		h = sha256.New()
	}
	_, err := io.WriteString(h, s)
	if err != nil {
		log.Println("Failed to write string value")
		os.Exit(1)
	}
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
