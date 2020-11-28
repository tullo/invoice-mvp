package rest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
	"github.com/tullo/invoice-mvp/identityprovider/fusionauth"
	"github.com/tullo/invoice-mvp/identityprovider/secret"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// Key is used to store/retrieve a Claims value from a context.Context.
const Key ctxKey = 1

// These are the expected values for Claims.Roles.
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

var (
	audience string
	issuer   string
	km       map[string]fusionauth.Key
	realm    string
)

// Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	Roles []string `json:"roles"`
	jwt.StandardClaims
}

// Valid is called during the parsing of a token.
func (c Claims) Valid(helper *jwt.ValidationHelper) error {
	for _, r := range c.Roles {
		switch r {
		case RoleAdmin, RoleUser: // Role is valid.
		default:
			return fmt.Errorf("invalid role %q", r)
		}
	}
	if err := c.StandardClaims.Valid(helper); err != nil {
		return errors.Wrap(err, "validating standard claims")
	}
	return nil
}

// Authorized returns true if claims has at least one of the provided roles.
func (c Claims) Authorized(roles ...string) bool {
	for _, has := range c.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}
	return false
}

// Realm returns the configured realm by consulting the
// environment variable "AUTH_REALM".
func Realm() string {
	if len(realm) > 0 {
		return realm
	}

	if v, ok := os.LookupEnv("AUTH_REALM"); ok {
		realm = v
	}

	return realm
}

// IDPAudience returns the configured issuer by consulting the
// environment variable "IDP_ISSUER".
func IDPAudience() string {
	if len(audience) > 0 {
		return audience
	}

	if v, ok := os.LookupEnv("CLIENT_ID"); ok {
		audience = v
	}

	return audience
}

// IDPIssuer returns the configured issuer by consulting the
// environment variable "IDP_ISSUER".
func IDPIssuer() string {
	if len(issuer) > 0 {
		return issuer
	}

	if v, ok := os.LookupEnv("IDP_ISSUER"); ok {
		issuer = v
	}

	return issuer
}

// Decorator func
func decorator(f func()) func() {
	return func() {
		log.Println("before fn call")
		f()
		log.Println("after fn call")
	}
}

// ===== BASIC AUTH ===========================================================

// BasicAuth decorator
func BasicAuth(next Handler) Handler {
	// closure func block
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			if username == os.Getenv("MVP_USERNAME") && password == os.Getenv("MVP_PASSWORD") {
				next(ctx, w, r) // call request handler
				return
			}
		}
		w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", realm))
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// ===== DIGEST AUTH ==========================================================

const password = "time"
const nonce = "UAZs1dp3wX5BtXEpoCXKO2lHhap564rX"
const opaque = "XF3tAJ3483jUUAUJJQJJAHDQP01MJHD"
const digest = "Digest"

// DigestAuth decorator
func DigestAuth(next Handler) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, digest) {
			m := digestParts(auth[len(digest)+1:])
			// calculate response hash (h3) using shared secret: password
			h1 := hash(fmt.Sprintf("%s:%s:%s", m["username"], realm, password))
			h2 := hash(fmt.Sprintf("%s:%s", r.Method, m["uri"]))
			nc, err := strconv.ParseInt(m["nc"], 16, 64)
			if err != nil {
				log.Println("Failed to parse hex value:", m["nc"])
				os.Exit(1)
			}
			h3 := hash(fmt.Sprintf("%s:%s:%08x:%s:%s:%s", h1, m["nonce"], nc, m["cnonce"], m["qop"], h2))

			if h3 == m["response"] { // client and server response hashes must match
				next(ctx, w, r) // call request handler
				return
			}
		}
		a := fmt.Sprintf(`Digest realm=%q, nonce=%q, opaque=%q, qop="auth", algorithm="SHA-256"`, realm, nonce, opaque)
		w.Header().Set("WWW-Authenticate", a)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func hash(text string) string {
	h := sha256.New()
	_, err := io.WriteString(h, text)
	if err != nil {
		log.Println("Failed to write string value")
		os.Exit(1)
	}
	hash := h.Sum(nil)
	// log.Printf("hash: %x %s\n", hash, text)
	return hex.EncodeToString(hash)
}

func digestParts(s string) map[string]string {
	m := map[string]string{}
	parts := []string{"username", "nonce", "realm", "qop", "uri", "nc", "response", "opaque", "cnonce"}
	for _, p := range strings.Split(s, ",") {
		for _, part := range parts {
			if strings.Contains(p, part) {
				s := strings.Split(strings.TrimSpace(p), "=")
				if s[0] == part {
					m[part] = strings.Trim(s[1], `"`)
				}
			}
		}
	}
	return m
}

// ===== JWT AUTH =============================================================

// JWTAuth decorator
func JWTAuth(next Handler) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token := ExtractJwt(r.Header)
		if claims, ok := verifyJWT(token); ok {
			// Add claims to the context so they can be retrieved later.
			ctx = context.WithValue(ctx, Key, claims)
			next(ctx, w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", realm))
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// ExtractJwt extracts the jwt token from the header line.
func ExtractJwt(h http.Header) string {
	var jwtRegex = regexp.MustCompile(`^Bearer (\S+)$`)

	if hs, ok := h["Authorization"]; ok {
		for _, value := range hs {
			if ss := jwtRegex.FindStringSubmatch(value); ss != nil {
				return ss[1]
			}
		}
	}

	return ""
}

// HMACKeyFunc verifies the token signing method and returns the shared HMAC
// secret as key used for signature validation.
func HMACKeyFunc(t *jwt.Token) (interface{}, error) {
	// signing method from token header must match expected method.
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}
	// Public IDP Key = shared HMAC secret
	return []byte(secret.Shared), nil
}

func verifyJWT(s string) (Claims, bool) {
	var claims Claims
	if len(s) == 0 {
		log.Println("Got an empty token: unable to verify")
		return claims, false
	}

	// Verify unparsed token parts.
	parts := strings.Split(s, ".")
	if len(parts) < 3 {
		log.Println("Invalid token parts count: unable to verify")
		return claims, false
	}

	var po []jwt.ParserOption
	po = append(po, jwt.WithIssuer(IDPIssuer()))
	po = append(po, jwt.WithAudience(IDPAudience()))
	t, err := jwt.ParseWithClaims(s, &claims, RS256KeyFunc, po...)
	if err != nil {
		log.Println("Token parsing:", err)
		return claims, false
	}

	if !t.Valid {
		log.Println("token is not valid")
		return claims, false
	}

	return claims, t.Valid
}

// OAuth2AccessCodeGrant decorator makes sure the redirect URI is valid.
func OAuth2AccessCodeGrant(next Handler) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		var code string
		var state string
		if v, ok := q["code"]; ok {
			code = v[0]
		}
		if v, ok := q["userState"]; ok {
			state = v[0]
		}

		if len(code) > 0 && len(state) > 0 && state == "Authenticated" {
			next(ctx, w, r) // call request handler
			return
		}

		w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", realm))
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

// OAuth2AccessTokenHandler exchanges the oauth code grant for an access token.
func (a Adapter) OAuth2AccessTokenHandler() Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		//dump, _ := httputil.DumpRequest(r, true)
		//fmt.Println("request:", string(dump))
		q := r.URL.Query()
		var code string
		if v, ok := q["code"]; ok {
			code = v[0]
		}

		t, err := a.exchangeOAuthCodeForAccessToken(code)
		if err != nil {
			log.Println(err)
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", realm))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if b, err := json.Marshal(&t); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (a Adapter) exchangeOAuthCodeForAccessToken(codeGrant string) (fusionauth.AuthInfo, error) {

	// Form data
	data := url.Values{}
	// data.Set("user_code", "")
	// data.Set("scope", "")
	data.Set("code", codeGrant)
	data.Set("client_id", a.idp.ClientID)
	data.Set("client_secret", a.idp.ClientSecret)
	data.Set("grant_type", a.idp.GrantType)
	data.Set("redirect_uri", a.idp.RedirectURI)

	var ai fusionauth.AuthInfo
	var c http.Client
	res, err := c.PostForm(a.idp.TokenURI, data)
	if err != nil {
		return ai, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ai, err
	}
	if err := json.Unmarshal(body, &ai); err != nil {
		return ai, err
	}
	return ai, nil
}

// ===== IDENTITY PROVIDER ====================================================

// RS256KeyFunc verifies the token signing method and returns
// the public signing key, matching the key identified by the
// tokens "kid" header value, used for signature validation.
func RS256KeyFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Header["alg"]; !ok {
		return nil, fmt.Errorf("Expected 'alg' to exist in token header")
	}
	if _, ok := t.Header["kid"]; !ok {
		return nil, fmt.Errorf("Expected 'kid' to exist in token header")
	}
	halg := stringVal(t.Header["alg"])
	if len(halg) < 1 {
		return nil, fmt.Errorf("Unexpected 'alg' value,  got: %v", t.Header["alg"])
	}
	hkid := stringVal(t.Header["kid"])
	if len(hkid) < 1 {
		return nil, fmt.Errorf("Unexpected 'kid' value,  got: %v", t.Header["kid"])
	}

	// use the prefetched public signing key
	key, ok := km[hkid]
	if !ok {
		log.Printf("JSON key with ID=[%s] not found.", hkid)
		// This execution path covers three cases:
		if len(km) < 1 {
			// (1) The app was just launched and no keys are fetched from the IDP yet.
			log.Print("Loading JSON Key Set from IDP service.")
			cnt, err := loadJSONKeySet()
			if err != nil {
				log.Println("loadJSONKeySet", err.Error())
			}
			log.Printf("Loaded keyset contains [%d] keys.\n", cnt)
		} else {
			// (2) A new signing key has been setup for this app in the IDP configuration.
			log.Print("Reloading JSON Key Set from IDP service")
			cnt, err := loadJSONKeySet()
			if err != nil {
				log.Println("loadJSONKeySet", err.Error())
			}
			log.Printf("Reloaded keyset contains [%d] keys.\n", cnt)
		}

		key, ok = km[hkid]
		if !ok {
			// (3) The kid does not match any of the published signing keys.
			return nil, fmt.Errorf("Key ID %s not found in published IDP key set", hkid)
		}
	}

	// Check signing method in use, Only RS256 is supported.
	m := jwt.GetSigningMethod(key.Alg)
	if m.Alg() != halg || m.Alg() != "RS256" {
		return nil, fmt.Errorf("Unexpected signing method in token-header: %v, expected: %v", halg, m.Alg())
	}

	// Public Key instance of the Token-Issuer
	return key.Instance, nil
}

func loadJSONKeySet() (int, error) {
	// Retrieve the published JSON Web Key Set (JWKS).
	ks, err := fusionauth.JSONWebKeySet(fusionauth.JWKSEndpoint)
	if err != nil {
		log.Fatal("Could not retrieve published key set")
	}
	useFilter := "sig"
	km = fusionauth.PublicSigningKeyMap(ks.Keys, useFilter)
	km, err = fusionauth.RetrievePublicKeyInstances(km)
	if err != nil {
		return 0, fmt.Errorf("loadJSONKeySet() expected no error, got [%v]", err)
	}
	return len(km), nil
}

// Pulls out the concrete string value of the interface.
func stringVal(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	default:
		return ""
	}
}
