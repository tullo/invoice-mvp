package fusionauth

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
)

const (
	// AuthorizeEndpoint ...
	AuthorizeEndpoint = "http://localhost:9011/oauth2/authorize"
	// JWKSEndpoint ...
	JWKSEndpoint = "http://localhost:9011/.well-known/jwks.json"
	// LoginEndpoint ...
	LoginEndpoint = "http://localhost:9011/api/login"
	// PublicKeyEndpoint ...
	PublicKeyEndpoint = "http://localhost:9011/api/jwt/public-key"
	// TokenEndpoint ...
	TokenEndpoint = "http://localhost:9011/oauth2/token"
)

// Key represents a JSON Web Key
type Key struct {
	Alg          string `json:"alg"`
	ID           string `json:"kid"`
	PublicKeyPEM string `json:"publicKey"`
	Use          string `json:"use"`
	Instance     *rsa.PublicKey
}

// KeySet holds a JSON Web Key Set (JWKS)
type KeySet struct {
	Keys []Key `json:"keys,omitempty"`
}

// AuthInfo represents incomming data from the identity provider.
type AuthInfo struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
	TokenType   string  `json:"token_type"`
	UserID      string  `json:"userId"`
}

// AuthConfig holds configuration for external IDP integration.
type AuthConfig struct {
	ClientID     string
	ClientSecret string
	GrantType    string
	Issuer       string
	TokenURI     string
	RedirectURI  string // Must match FA config for "Authorized redirect URLs"
}

func tlsTransportConfig() (http.RoundTripper, error) {

	// Read the CA signed certificate.
	cert, err := ioutil.ReadFile(filepath.Join(os.Getenv("BASE_DIR"), "localhost+2.pem"))
	if err != nil {
		return nil, errors.Wrap(err, "loading certificate file")
	}

	// Add cert to the pool of trusted certificats.
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		return nil, errors.Wrap(err, "appending certificate to cert pool")
	}

	// Use trusted server certificate for client transport config.
	var c tls.Config
	c.RootCAs = certPool
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.TLSClientConfig = &c

	return t, nil
}

// Client returns a http.Client instance with TLS transport configured if tls
// is set to true.
func Client(tls bool) (*http.Client, error) {
	var c http.Client
	if tls {
		t, err := tlsTransportConfig()
		if err != nil {
			log.Println("(tls) err:", err)
			return &c, errors.Wrap(err, "creating transport TLS client config")
		}
		// Client makes sure to only talk to servers in possession
		// of the private key used to sign the certificate.
		c.Transport = t
		return &c, nil
	}
	c.Transport = http.DefaultTransport.(*http.Transport).Clone()
	return &c, nil
}

func accessCodeGrant(data url.Values) (string, error) {
	var code string

	form := make(url.Values)
	form.Set("response_type", "code")
	form.Set("client_id", os.Getenv("CLIENT_ID"))       // Invoice MVP
	form.Set("redirect_uri", os.Getenv("REDIRECT_URI")) // Must match FA config.

	// Overwrite default settings
	if ss, ok := data["client_id"]; ok {
		form.Set("client_id", ss[0])
	}
	// Credentials
	if ss, ok := data["loginId"]; ok {
		form.Set("loginId", ss[0])
	}
	if ss, ok := data["password"]; ok {
		form.Set("password", ss[0])
	}

	// URL-encoded payload
	payload := form.Encode()
	req, err := http.NewRequest("POST", AuthorizeEndpoint, strings.NewReader(payload))
	if err != nil {
		return code, errors.Wrap(err, "posting URL-encoded payload")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	// HTTP roundtrip, redirects are not followed.
	var noRedirect http.RoundTripper = &http.Transport{}
	res, err := noRedirect.RoundTrip(req)
	defer res.Body.Close()
	if err != nil {
		return code, errors.Wrap(err, "posting URL-encoded payload")
	}

	if res.StatusCode != http.StatusFound {
		return code, errors.Errorf("Unexpected response status %d", res.StatusCode)
	}
	if _, ok := res.Header["Location"]; !ok {
		return code, errors.New("Location header not available")
	}
	if loc, ok := res.Header["Location"]; ok {
		if len(loc) < 1 {
			return code, errors.New("Location value missing")
		}
	}

	loc := res.Header["Location"][0]
	u, err := url.Parse(loc)
	if err != nil {
		return code, errors.Wrap(err, "parsing location header")
	}
	q := u.Query()
	if us, ok := q["userState"]; ok {
		if us[0] != "Authenticated" {
			return code, errors.Errorf("Unexpected user state %v", us)
		}
	}
	if _, ok := q["code"]; !ok {
		return code, errors.Errorf("Access Code Grant not found")
	}

	return q["code"][0], nil
}

// Login uses the provided user credentials to login with
// the IDM and converts the resulting code grant to JWT token.
func Login(data url.Values) (AuthInfo, error) {

	var auth AuthInfo

	grant, err := accessCodeGrant(data)
	if err != nil {
		return auth, errors.Wrap(err, "access code grant retrieval failed")
	}

	form := make(url.Values)
	//form.Set("user_code", "")
	//form.Set("scope", "")
	form.Set("code", grant)
	form.Set("grant_type", os.Getenv("GRANT_TYPE"))
	form.Set("redirect_uri", os.Getenv("REDIRECT_URI")) // Must match FA config.
	form.Set("client_id", os.Getenv("CLIENT_ID"))       // Invoice MVP
	form.Set("client_secret", os.Getenv("CLIENT_SECRET"))

	// Overwrite default settings
	if ss, ok := data["client_id"]; ok {
		form.Set("client_id", ss[0])
	}
	if ss, ok := data["client_secret"]; ok {
		form.Set("client_secret", ss[0])
	}

	client, err := Client(false) // no TLS
	if err != nil {
		return auth, errors.Wrap(err, "creating client instance")
	}
	payload := form.Encode() // URL-encoded payload
	res, err := client.Post(TokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(payload))
	if err != nil {
		return auth, errors.Wrap(err, "posting URL-encoded payload")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return auth, errors.Errorf("Unexpected response status %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return auth, errors.Wrap(err, "reading response body")
	}

	if err := json.Unmarshal(body, &auth); err != nil {
		return auth, errors.Wrap(err, "unmarshalling response body")
	}
	if auth.TokenType != "Bearer" {
		return auth, errors.Errorf("TokenType not valid %v", auth.TokenType)
	}
	if auth.ExpiresIn < 0 {
		return auth, errors.Errorf("Token expired %v", auth.ExpiresIn)
	}
	if len(auth.UserID) < 1 {
		return auth, errors.Errorf("UserID not valid %v", auth.UserID)
	}
	if len(auth.AccessToken) < 1 {
		return auth, errors.Errorf("AccessToken not valid %v", auth.AccessToken)
	}

	return auth, nil
}

// JSONWebKeySet retrieves the publisched set of JSON Web
// Keys from the identity provider.
func JSONWebKeySet(jwksURI string) (KeySet, error) {
	var ks KeySet
	client, err := Client(false) // no TLS
	if err != nil {
		return ks, errors.Wrap(err, "creating client instance")
	}
	res, err := client.Get(jwksURI)
	if err != nil {
		return ks, errors.Wrap(err, "retrieving json web key set")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ks, errors.Wrap(err, "reading response body")
	}

	if err := json.Unmarshal(body, &ks); err != nil {
		return ks, errors.Wrap(err, "unmarshalling response body")
	}
	return ks, nil
}

// PublicSigningKey retrieves the public signing key
// identified by the passed key ID.
func PublicSigningKey(keyID string) (Key, error) {
	var key Key
	keyURI := fmt.Sprintf("%s?kid=%s", PublicKeyEndpoint, keyID)
	client, err := Client(false) // no TLS
	if err != nil {
		return key, errors.Wrap(err, "creating client instance")
	}
	res, err := client.Get(keyURI)
	if err != nil {
		return key, errors.Wrap(err, "retrieving public signing key")
	}
	// dump, _ := httputil.DumpResponse(res, true)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return key, errors.Wrap(err, "reading response body")
	}

	if err := json.Unmarshal(body, &key); err != nil {
		return key, errors.Wrap(err, "unmarshalling response body")
	}

	return key, nil
}

// PublicSigningKeyMap filters the keyset by JWK property
// "use" and returns the resulting map.
func PublicSigningKeyMap(keys []Key, filter string) map[string]Key {
	var m = make(map[string]Key)
	for _, k := range keys {
		// JWK property `use` determines the JWK is for signature verification
		if k.Use == filter {
			m[k.ID] = Key{
				Alg: k.Alg,
				ID:  k.ID,
				Use: k.Use,
			}
		}
	}
	return m
}

// RetrievePublicKeyInstance gets the public signing key
// from the Identity Provider service and parses the PEM
// key representation into a key instance.
func RetrievePublicKeyInstance(keyID string) (Key, error) {
	var key Key
	k, err := PublicSigningKey(keyID)
	if err != nil {
		return key, errors.Wrap(err, "Could not retrieve public signing key from IDP")
	}
	key.PublicKeyPEM = k.PublicKeyPEM
	pkInstance, err := jwt.ParseRSAPublicKeyFromPEM([]byte(k.PublicKeyPEM))
	if err != nil {
		return key, errors.Wrap(err, "Could not parse public signing key")
	}
	key.Instance = pkInstance
	return key, nil
}

// RetrievePublicKeyInstances gets the public signing keys
// from the IDP and parses the PEM key representation.
func RetrievePublicKeyInstances(km map[string]Key) (map[string]Key, error) {
	im := make(map[string]Key, len(km))
	for _, v := range km {
		key, err := RetrievePublicKeyInstance(v.ID)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve public signing key instance")
		}
		// add existing key meta data
		key.Alg = v.Alg
		key.ID = v.ID
		key.Use = v.Use
		im[v.ID] = key
	}

	return im, nil
}

// ContainsValidSigningKey looks for the public signing key
// with (use=sign) and specified signing method (alg).
func ContainsValidSigningKey(ks []Key, alg string) bool {
	for _, k := range ks {
		// JWK property `use` determines the JWK is for signature verification
		if k.Alg == alg && k.Use == "sig" {
			return true
		}
	}
	return false
}
