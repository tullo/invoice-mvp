package roles

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/tullo/invoice-mvp/domain"
	"github.com/tullo/invoice-mvp/rest"
)

// These are the expected values for Claims.Roles.
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
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

// RoleRepository is a small interface used for assertions.
type RoleRepository interface {
	CustomerByID(id int) domain.Customer
	GetInvoice(id int, join ...string) domain.Invoice
}

// AssertAdmin decorator.
func AssertAdmin(next rest.Handler) rest.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token := rest.ExtractJwt(r.Header)
		if !isAdmin(token) {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", rest.Realm()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(ctx, w, r) // call request handler
	}
}

func isAdmin(s string) bool {
	// Parse and validate token using keyfunc.
	var claims Claims
	var po []jwt.ParserOption
	po = append(po, jwt.WithIssuer(rest.IDPIssuer()))
	po = append(po, jwt.WithAudience(rest.IDPAudience()))
	token, err := jwt.ParseWithClaims(s, &claims, rest.RS256KeyFunc, po...)
	if err != nil {
		log.Println("parsing token", err)
		return false
	}

	if !token.Valid {
		log.Println("token is not valid")
		return false
	}

	return claims.Authorized(RoleAdmin)
}

// AssertOwnsInvoice decorator
func AssertOwnsInvoice(next rest.Handler, rep RoleRepository) rest.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token := rest.ExtractJwt(r.Header)
		id, _ := strconv.Atoi(mux.Vars(r)["invoiceId"])
		// Load invoice
		i := rep.GetInvoice(id)
		// Load customer bound to invoice
		c := rep.CustomerByID(i.CustomerID)
		// Verify user owns the customer
		uid := rest.Claim(token, "sub")
		if c.UserID != uid {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next(ctx, w, r) // call request handler
	}
}
