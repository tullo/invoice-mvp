package roles

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/mux"
	"github.com/tullo/invoice-mvp/domain"
	"github.com/tullo/invoice-mvp/rest"
)

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
	t, err := jwt.Parse(s, rest.HMACKeyFunc)
	if err != nil {
		return false
	}

	// Extract admin claim.
	if cm, ok := t.Claims.(jwt.MapClaims); ok {
		if cm["admin"] != nil {
			return cm["admin"].(bool) // map[string]interface{}
		}
	}

	return false
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
