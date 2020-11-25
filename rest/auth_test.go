package rest_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestAuthAccessCodeGrant(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
		t.Error(err)
	}

	// POST form data
	data := url.Values{}
	data.Set("response_type", "code")
	data.Set("loginId", os.Getenv("TEST_LOGIN"))
	data.Set("password", os.Getenv("TEST_PASSWD"))
	data.Set("client_id", os.Getenv("CLIENT_ID"))                 // Invoice MVP
	data.Set("redirect_uri", "https://127.0.0.1:8443/authorized") // Must match FA config.

	// FA authorize:
	endpoint := "http://localhost:9011/oauth2/authorize"
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		t.Error(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	var noRedirect http.RoundTripper = &http.Transport{}
	res, err := noRedirect.RoundTrip(r)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusFound {
		t.Error("Unexpected response status", res.StatusCode)
	}
	if _, ok := res.Header["Location"]; !ok {
		t.Error("Location header not available")
	}
	if loc, ok := res.Header["Location"]; ok {
		if len(loc) < 1 {
			t.Error("Location value missing")
		}
	}
	loc := res.Header["Location"][0]
	u, err := url.Parse(loc)
	if err != nil {
		t.Error(err)
	}
	q := u.Query()
	if us, ok := q["userState"]; ok {
		if us[0] != "Authenticated" {
			t.Error("Unexpected user state", us)
		}
	}
	t.Log(q["code"][0])
}

func TestAuthAccessToken(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
		t.Error(err)
	}

	// ====================== AUTHORIZE =======================================
	// POST form data
	data := url.Values{}
	data.Set("response_type", "code")
	data.Set("loginId", os.Getenv("TEST_LOGIN"))
	data.Set("password", os.Getenv("TEST_PASSWD"))
	data.Set("client_id", os.Getenv("CLIENT_ID"))                 // Invoice MVP
	data.Set("redirect_uri", "https://127.0.0.1:8443/authorized") // must match FA config.

	// FA authorize:
	endpoint := "http://localhost:9011/oauth2/authorize"
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		t.Error(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	var noRedirect http.RoundTripper = &http.Transport{}
	res, err := noRedirect.RoundTrip(r)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusFound {
		t.Error("Unexpected response status", res.StatusCode)
	}
	if _, ok := res.Header["Location"]; !ok {
		t.Error("Location header not available")
	}
	if loc, ok := res.Header["Location"]; ok {
		if len(loc) < 1 {
			t.Error("Location value missing")
		}
	}
	loc := res.Header["Location"][0]
	u, err := url.Parse(loc)
	if err != nil {
		t.Error(err)
	}
	q := u.Query()
	if us, ok := q["userState"]; ok {
		if us[0] != "Authenticated" {
			t.Error("Unexpected user state", us)
		}
	}

	// ====================== ACCESS TOKEN ====================================

	// FusionAuth token endpoint:
	endpoint = "http://localhost:9011/oauth2/token"
	data = url.Values{}
	//data.Set("user_code", "")
	//data.Set("scope", "")
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", "https://127.0.0.1:8443/authorized") // Must match FA config.
	data.Set("client_id", os.Getenv("CLIENT_ID"))                 // Invoice MVP
	data.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	data.Set("code", q["code"][0])

	var client http.Client
	r, err = http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		t.Error(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err = client.Do(r)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Error("Unexpected response status", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	var auth struct {
		AccessToken string  `json:"access_token"`
		ExpiresIn   float64 `json:"expires_in"`
		TokenType   string  `json:"token_type"`
		UserID      string  `json:"userId"`
	}
	if err := json.Unmarshal(body, &auth); err != nil {
		t.Fatal(err)
	}
	if auth.TokenType != "Bearer" {
		t.Error("TokenType not valid", auth.TokenType)
	}
	if auth.ExpiresIn < 0 {
		t.Error("Token expired", auth.ExpiresIn)
	}
	if len(auth.UserID) < 1 {
		t.Error("UserID not valid", auth.UserID)
	}
	if len(auth.AccessToken) < 1 {
		t.Error("AccessToken not valid", auth.AccessToken)
	}
}
