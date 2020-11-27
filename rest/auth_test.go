package rest_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/tullo/invoice-mvp/identityprovider/fusionauth"
	"github.com/tullo/invoice-mvp/rest"
	"gopkg.in/go-playground/assert.v1"
)

func TestAuthAccessCodeGrant(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
		t.Fatal(err)
	}

	// POST form data
	data := url.Values{}
	data.Set("response_type", "code")
	data.Set("loginId", os.Getenv("TEST_LOGIN"))
	data.Set("password", os.Getenv("TEST_PASSWD"))
	data.Set("client_id", os.Getenv("TEST_CLIENT_ID"))  // Test application
	data.Set("redirect_uri", os.Getenv("REDIRECT_URI")) // Must match FA config.

	// FA authorize:
	endpoint := "http://localhost:9011/oauth2/authorize"
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	var noRedirect http.RoundTripper = &http.Transport{}
	res, err := noRedirect.RoundTrip(r)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusFound {
		t.Fatal("Unexpected response status", res.StatusCode)
	}
	if _, ok := res.Header["Location"]; !ok {
		t.Fatal("Location header not available")
	}
	if loc, ok := res.Header["Location"]; ok {
		if len(loc) < 1 {
			t.Fatal("Location value missing")
		}
	}
	loc := res.Header["Location"][0]
	u, err := url.Parse(loc)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	if us, ok := q["userState"]; ok {
		if us[0] != "Authenticated" {
			t.Error("Unexpected user state", us)
		}
	}
	// t.Log(q["code"][0])
}

func TestAuthAccessToken(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
		t.Fatal(err)
	}

	// ====================== AUTHORIZE =======================================
	// POST form data
	data := url.Values{}
	data.Set("response_type", "code")
	data.Set("loginId", os.Getenv("TEST_LOGIN"))
	data.Set("password", os.Getenv("TEST_PASSWD"))
	data.Set("client_id", os.Getenv("TEST_CLIENT_ID"))  // Test application
	data.Set("redirect_uri", os.Getenv("REDIRECT_URI")) // Must match FA config.

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
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusFound {
		t.Fatal("Unexpected response status", res.StatusCode)
	}
	if _, ok := res.Header["Location"]; !ok {
		t.Fatal("Location header not available")
	}
	if loc, ok := res.Header["Location"]; ok {
		if len(loc) < 1 {
			t.Fatal("Location value missing")
		}
	}
	loc := res.Header["Location"][0]
	u, err := url.Parse(loc)
	if err != nil {
		t.Fatal(err)
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
	data.Set("redirect_uri", os.Getenv("REDIRECT_URI")) // Must match FA config.
	data.Set("client_id", os.Getenv("TEST_CLIENT_ID"))  // Test application
	data.Set("client_secret", os.Getenv("TEST_CLIENT_SECRET"))
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

func TestOAuth2Handler(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
		t.Fatal(err)
	}

	// ====================== AUTHORIZE =======================================
	// POST form data
	data := url.Values{}
	data.Set("response_type", "code")
	data.Set("loginId", os.Getenv("MVP_USERNAME"))
	data.Set("password", os.Getenv("MVP_PASSWORD"))
	data.Set("client_id", os.Getenv("CLIENT_ID"))       // Test application
	data.Set("redirect_uri", os.Getenv("REDIRECT_URI")) // Must match FA config.

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
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusFound {
		t.Error("Unexpected response status", res.StatusCode)
	}
	if _, ok := res.Header["Location"]; !ok {
		t.Fatal("Location header not available")
	}
	if loc, ok := res.Header["Location"]; ok {
		if len(loc) < 1 {
			t.Fatal("Location value missing")
		}
	}
	loc := res.Header["Location"][0]
	//t.Log("Location", loc)
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

	//=========================================================================
	// Exchange access code grant with JWT access token

	rr := httptest.NewRecorder()
	a := rest.NewAdapter()
	// IDP redirects to this URI after user authentication
	a.Handle("/auth/token", rest.OAuth2AccessCodeGrant(a.OAuth2AccessTokenHandler())).Methods("GET")

	req, _ := http.NewRequest("GET", loc, nil)
	a.R.ServeHTTP(rr, req)

	//=========================================================================
	// Assert

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	//dump, _ := httputil.DumpResponse(rr.Result(), true)
	//t.Log("responseDump:\n", string(dump))
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Error(err)
	}

	var auth fusionauth.AuthInfo
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
