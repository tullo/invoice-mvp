package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/tullo/invoice-mvp/domain"
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
	_, status := activities(client, jwtToken)
	log.Println("Initial response:", status)
	if status == http.StatusUnauthorized {
		log.Println("Login and retrieve JWT access token from IDP.")
		data := make(url.Values)
		data.Set("loginId", os.Getenv("MVP_USERNAME"))
		data.Set("password", os.Getenv("MVP_PASSWORD"))
		auth, err := fusionauth.Login(data)
		if err != nil {
			log.Fatal(err)
		}
		as, status := activities(client, auth.AccessToken)
		if status != http.StatusOK {
			log.Fatal("Authorized response:", status)
		}
		log.Printf("Activities managed by user %s:\n", auth.UserID)
		for _, a := range as {
			log.Printf("Activity: %+q\n", a)
		}
	}
}

func activities(client *http.Client, token string) ([]domain.Activity, int) {
	var as []domain.Activity
	req, _ := http.NewRequest("GET", "https://127.0.0.1:8443/activities", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		log.Printf("failed to retrieve activities: %v", err)
		return as, res.StatusCode
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return as, res.StatusCode
	}
	if err := json.Unmarshal(body, &as); err != nil {
		log.Println(err)
		return as, res.StatusCode
	}
	return as, http.StatusOK
}
