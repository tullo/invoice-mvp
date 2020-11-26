package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/tullo/invoice-mvp/identityprovider/fusionauth"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := fusionauth.Client(true)
	if err != nil {
		log.Fatal(err)
	}

	jwtToken := "header.payload.signature"
	status := activities(client, jwtToken)
	log.Println("Activities: initial response:", status)
	if status == http.StatusUnauthorized {
		log.Println("Login and retrieve JWT access token.")
		data := make(url.Values)
		data.Set("loginId", os.Getenv("TEST_LOGIN"))
		data.Set("password", os.Getenv("TEST_PASSWD"))
		auth, err := fusionauth.Login(data)
		if err != nil {
			log.Fatal(err)
		}
		jwtToken = auth.AccessToken
		status = activities(client, jwtToken)
		log.Println("Activities: Authorized Request:", status)
	}
}

func activities(client *http.Client, token string) int {
	log.Println("Activities: token:", token)

	req, _ := http.NewRequest("GET", "https://127.0.0.1:8443/activities", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	res.Body.Close()
	if err != nil {
		log.Printf("%v", err)
		return http.StatusInternalServerError
	}
	return res.StatusCode
}
