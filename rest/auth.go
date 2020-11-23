package rest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/tullo/invoice-mvp/identityprovider/secret"
)

// Decorator func
func decorator(f func()) func() {
	return func() {
		log.Println("before fn call")
		f()
		log.Println("after fn call")
	}
}

// BasicAuth decorator
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	// closure func block
	return func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			if username == os.Getenv("MVP_USERNAME") && password == os.Getenv("MVP_PASSWORD") {
				next.ServeHTTP(w, r) // call request handler
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="invoice.mvp"`)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// ===== DIGEST AUTH ==========================================================

const password = "time"
const realm = "invoice.mvp"
const nonce = "UAZs1dp3wX5BtXEpoCXKO2lHhap564rX"
const opaque = "XF3tAJ3483jUUAUJJQJJAHDQP01MJHD"
const digest = "Digest"

// DigestAuth decorator
func DigestAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
				next.ServeHTTP(w, r) // call request handler
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
func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractJwt(r.Header)
		if verifyJWT(token) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", `Bearer realm="invoice.mvp"`)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func extractJwt(h http.Header) string {
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

func verifyJWT(s string) bool {
	// Parse and validate token using a keyfunc.
	t, err := jwt.Parse(s, func(t *jwt.Token) (interface{}, error) {
		// signing method from token header must match expected method.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		// Public IDP Key = shared HMAC secret
		return []byte(secret.Shared), nil
	})

	return err == nil && t.Valid
}
