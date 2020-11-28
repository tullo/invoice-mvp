package roles

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		if !isAdmin(ctx) {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", rest.Realm()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(ctx, w, r) // call request handler
	}
}

func isAdmin(ctx context.Context) bool {
	claims, ok := ctx.Value(rest.Key).(rest.Claims)
	if !ok {
		log.Println("claims missing from context")
		return false
	}

	return claims.Authorized(rest.RoleAdmin)
}

// AssertOwnsInvoice decorator
func AssertOwnsInvoice(next rest.Handler, rep RoleRepository) rest.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["invoiceId"])
		// Load invoice
		i := rep.GetInvoice(id)
		// Load customer bound to invoice
		c := rep.CustomerByID(i.CustomerID)
		// Verify user owns the customer
		claims, ok := ctx.Value(rest.Key).(rest.Claims)
		if !ok {
			log.Println("claims missing from context")
			return
		}
		uid := claims.Subject
		if c.UserID != uid {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next(ctx, w, r) // call request handler
	}
}
