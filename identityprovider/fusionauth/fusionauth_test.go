package fusionauth_test

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/joho/godotenv"
	"github.com/tullo/invoice-mvp/identityprovider/fusionauth"
)

func Test_client(t *testing.T) {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error loading .env file")
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		wantTLS bool
		wantErr bool
	}{
		{
			"successful test with TLS",
			true,
			false,
		},

		{
			"successful test without TLS",
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := fusionauth.Client(tt.wantTLS)
			if (err != nil) != tt.wantErr {
				t.Errorf("client() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			// client should not be set
			if c == nil {
				t.Errorf("client() client: %v", c)
				return
			}
			// check our tls config
			if tt.wantTLS {
				// should have configured custom transport
				if _, ok := c.Transport.(*http.Transport); !ok {
					t.Error("client() Transport expected, got none")
					return
				}
				// CA pool should contain a root cert
				tr := c.Transport.(*http.Transport)
				rCAs := tr.TLSClientConfig.RootCAs
				if len(rCAs.Subjects()) < 1 {
					t.Error("client() RootCAs expected, got none")
					return
				}
				// subject should match
				var subject pkix.RDNSequence
				sub := "OU=anda@pop-os (Andreas),O=mkcert development certificate"
				rawSubj := rCAs.Subjects()[0]
				if _, err := asn1.Unmarshal(rawSubj, &subject); err != nil {
					t.Errorf("client() Unmarshal err: %v", err)
					return
				}
				if sub != subject.String() {
					t.Errorf("client() client: subject: %v", subject)
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	// skips this test if short flag was provided.
	if testing.Short() {
		t.Skip()
	}
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error loading .env file")
		t.Fatal(err)
	}

	type args struct {
		loginId       string
		password      string
		client_id     string
		client_secret string
		redirect_uri  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantAlg string
	}{
		{
			name: "Valid Login, App Specific Token Signing Method",
			args: args{
				os.Getenv("TEST_LOGIN"),
				os.Getenv("TEST_PASSWD"),
				os.Getenv("CLIENT_ID"),
				os.Getenv("CLIENT_SECRET"),
				os.Getenv("REDIRECT_URI"),
			},
			wantErr: false,
			wantAlg: "RS256",
		},
		{
			name: "Default Token Signing Method",
			args: args{
				os.Getenv("TEST_LOGIN"),
				os.Getenv("TEST_PASSWD"),
				os.Getenv("TEST_CLIENT_ID"),     // Test app
				os.Getenv("TEST_CLIENT_SECRET"), // Test app secret
				os.Getenv("REDIRECT_URI"),
			},
			wantErr: false,
			wantAlg: "HS256",
		},
		{
			name: "Wrong Password",
			args: args{
				os.Getenv("TEST_LOGIN"),
				"T0ps3cr3t",
				os.Getenv("CLIENT_ID"),
				os.Getenv("CLIENT_SECRET"),
				os.Getenv("REDIRECT_URI"),
			},
			wantErr: true,
			wantAlg: "RS256",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// POST form data
			data := make(url.Values)
			data.Set("loginId", tt.args.loginId)
			data.Set("password", tt.args.password)
			data.Set("client_id", tt.args.client_id)
			data.Set("client_secret", tt.args.client_secret)
			data.Set("redirect_uri", tt.args.redirect_uri)

			auth, err := fusionauth.Login(data)
			if !tt.wantErr && err != nil {
				t.Errorf("Login() got error: [%v], but did not expect one: %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err == nil {
				t.Errorf("Login() no error: [%v], but expected one: %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				//t.Logf("Login() got error: [%v], and expected one: %v", err, tt.wantErr)
				return
			}

			var claims jwt.Claims
			token, _, _ := jwt.NewParser().ParseUnverified(auth.AccessToken, claims)
			//for k, v := range token.Header {
			//	t.Logf("claim: %s=%v", k, v)
			//}
			if _, ok := token.Header["alg"].(string); !ok {
				t.Error("Login() token alg header not found")
				return
			}
			alg := token.Header["alg"].(string)
			if alg != tt.wantAlg {
				t.Errorf("Login() unexpected token alg header: got %v want %v", alg, tt.wantAlg)
			}
		})
	}
}

func TestJSONWebKeySet(t *testing.T) {
	type args struct {
		jwksURI string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Load JSON Web Keys from Identity Provider.",
			args:    args{fusionauth.JWKSEndpoint},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks, err := fusionauth.JSONWebKeySet(tt.args.jwksURI)
			if !tt.wantErr && err != nil {
				t.Errorf("JSONWebKeySet() got error: [%v], but did not expect one: %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err == nil {
				t.Errorf("JSONWebKeySet() no error: [%v], but expected one: %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				t.Logf("JSONWebKeySet() got error: [%v], and expected one: %v", err, tt.wantErr)
				return
			}

			if len(ks.Keys) < 1 {
				t.Errorf("JSONWebKeySet() expected json keyset: %v", ks.Keys)
				return
			}

			for _, k := range ks.Keys {
				if len(k.Alg) < 1 || len(k.Use) < 1 || len(k.ID) < 1 {
					t.Errorf("Unexpected key properties: %+v", k)
				}
			}
			//t.Errorf("JSONWebKeySet() keyset: %+v", ks)
		})
	}
}

func TestPublicSigningKey(t *testing.T) {
	type args struct {
		keyURI string
	}
	tests := []struct {
		name    string
		alg     string
		use     string
		wantErr bool
	}{
		{
			name:    "JSON Web Keys from IDP",
			alg:     "RS256",
			use:     "sig",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks, err := fusionauth.JSONWebKeySet(fusionauth.JWKSEndpoint)
			if err != nil {
				t.Error(err)
				return
			}
			// keyset should only contain public signing keys (use=sign)
			// with specified signing method (alg)
			if ok := fusionauth.ContainsValidSigningKey(ks.Keys, tt.alg); !ok {
				t.Errorf("PublicSigningKey() no public signing keys in json key set [%v].", ks.Keys)
			}
			m := fusionauth.PublicSigningKeyMap(ks.Keys, tt.use)
			for _, value := range m {
				k, err := fusionauth.PublicSigningKey(value.ID)
				if !tt.wantErr && err != nil {
					t.Errorf("PublicSigningKey() got error: [%v], but did not expect one: %v", err, tt.wantErr)
					return
				}
				if tt.wantErr && err == nil {
					t.Errorf("PublicSigningKey() no error: [%v], but expected one: %v", err, tt.wantErr)
					return
				}
				if tt.wantErr && err != nil {
					t.Errorf("PublicSigningKey() got error: [%v], and expected one: %v", err, tt.wantErr)
					return
				}
				if !strings.HasPrefix(k.PublicKeyPEM, "-----BEGIN PUBLIC KEY-----") {
					t.Errorf("JSONWebKeySet() expected pub sign key: %v", k.PublicKeyPEM)
					return
				}
			}
		})
	}
}

func TestPublicSigningKeyMap(t *testing.T) {

	ks, err := fusionauth.JSONWebKeySet(fusionauth.JWKSEndpoint)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name   string
		filter string
		items  int
	}{
		{
			name:   "Filter JWKs for signature verification",
			filter: "sig",
			items:  1,
		},
		{
			name:   "Empty map",
			filter: "xyz",
			items:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := fusionauth.PublicSigningKeyMap(ks.Keys, tt.filter)
			if len(m) < tt.items {
				t.Errorf("PublicSigningKeyMap() expected map with [%d] items, m=%v", tt.items, m)
				return
			}
		})
	}
}

func TestRetrievePublicKeyInstances(t *testing.T) {

	tests := []struct {
		name   string
		filter string
		items  int
	}{
		{
			name:   "Filter JWKs for signature verification",
			filter: "sig",
			items:  1,
		},
	}

	// Load published JSON key set.
	ks, err := fusionauth.JSONWebKeySet(fusionauth.JWKSEndpoint)
	if err != nil {
		t.Error(err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Filter public signing keys
			km := fusionauth.PublicSigningKeyMap(ks.Keys, tt.filter)
			if len(km) < tt.items {
				t.Errorf("PublicSigningKeyMap() expected map with [%d] items, m=%v", tt.items, km)
				return
			}
			im, err := fusionauth.RetrievePublicKeyInstances(km)
			if err != nil {
				t.Errorf("RetrievePublicKeyInstances() expected no error, got [%v]", err)
			}

			for _, k := range im {
				if k.Instance == nil {
					t.Error("expected non-nil key instance")
					return
				}
				if len(k.PublicKeyPEM) < 1 {
					t.Error("expected valid PEM rep of the key")
					return
				}
			}

		})
	}
}
