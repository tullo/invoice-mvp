package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tullo/invoice-mvp/identityprovider/fusionauth"

	"github.com/gorilla/mux"
	"github.com/tullo/invoice-mvp/domain"
	"github.com/tullo/invoice-mvp/usecase"
)

const dateFormat = "Mon, _2 Jan 2006 15:04:05 GMT"

func truncateToSeconds(t time.Time) time.Time {
	return t.Truncate(time.Duration(time.Second))
}

// Adapter converts HTTP request data into domain objects.
type Adapter struct {
	R   *mux.Router
	idp fusionauth.AuthConfig
}

// NewAdapter instantiates an adapter.
func NewAdapter() Adapter {
	var idp fusionauth.AuthConfig
	if v, ok := os.LookupEnv("CLIENT_ID"); ok {
		idp.ClientID = v
	}
	if v, ok := os.LookupEnv("CLIENT_SECRET"); ok {
		idp.ClientSecret = v
	}
	if v, ok := os.LookupEnv("GRANT_TYPE"); ok {
		idp.GrantType = v
	}
	if v, ok := os.LookupEnv("IDP_ISSUER"); ok {
		idp.Issuer = v
	}
	if v, ok := os.LookupEnv("REDIRECT_URI"); ok {
		idp.RedirectURI = v
	}
	if v, ok := os.LookupEnv("TOKEN_URI"); ok {
		idp.TokenURI = v
	}

	var a Adapter
	a.R = mux.NewRouter()
	a.idp = idp

	return a
}

// ListenAndServe launches a web server on port 8080.
func (a Adapter) ListenAndServe() {
	log.Printf("Listening on http://0.0.0.0%s\n", ":8080")
	_ = http.ListenAndServe(":8080", a.R)
}

// ListenAndServeTLS launches a web server on port 8080.
func (a Adapter) ListenAndServeTLS() {
	log.Printf("Listening on https://0.0.0.0%s\n", ":8443")
	_ = http.ListenAndServeTLS(":8443", "localhost+2.pem", "localhost+2-key.pem", a.R)
}

// Handler is a type that handles http requests.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)

// Handle creates a route and maps it to a path and handler.
func (a Adapter) Handle(path string, handler Handler) *mux.Route {

	h := func(w http.ResponseWriter, r *http.Request) {
		/*
			// Start or expand a distributed trace.
			ctx := r.Context()
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, r.URL.Path)
			defer span.End()

			// Set the context with the required values to
			// process the request.
			v := Values{
				TraceID: span.SpanContext().TraceID.String(),
				Now:     time.Now(),
			}
			ctx = context.WithValue(ctx, KeyValues, &v)
		*/
		handler(r.Context(), w, r)
	}
	return a.R.NewRoute().Path(path).HandlerFunc(h)
}

// InvoicePresenter returns a presenter matching the 'Accept' request header.
func (a Adapter) InvoicePresenter(w http.ResponseWriter, r *http.Request) (InvoicePresenter, bool) {
	// e.g. "Accept: application/json;q=0.8, application/hal+json"
	headers := strings.Split(r.Header.Get("Accept"), ",")
	var ip InvoicePresenter
	var ok bool
	for _, accept := range headers {
		switch accept {
		case "application/json", "application/hal+json":
			ip, ok = NewJSONInvoicePresenter(w), true
			break
		case "application/pdf":
			ip, ok = NewPDFInvoicePresenter(w, r), true
			break
		default:
			ip, ok = NewDefaultPresenter(), false
			break
		}
	}
	return ip, ok
}

// Extracts the authorized user's ID from the request (JWT).
func (a Adapter) currentUser(ctx context.Context) string {
	claims, ok := ctx.Value(Key).(Claims)
	if !ok {
		log.Println("claims missing from context")
		return ""
	}
	return claims.Subject
}

//=============================================================================
// Activity

func (a Adapter) readActivity(r *http.Request, uid string) (domain.Activity, error) {
	var act domain.Activity
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return act, err
	}
	if err := json.Unmarshal(body, &act); err != nil {
		return act, err
	}
	act.UserID = uid
	return act, nil
}

//=============================================================================
// Booking

func (a Adapter) readBooking(r *http.Request) (domain.Booking, error) {
	var b domain.Booking
	// extract invoiceId from the URI
	id, err := strconv.Atoi(mux.Vars(r)["invoiceId"])
	if err != nil {
		return b, err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return b, err
	}
	if err := json.Unmarshal(body, &b); err != nil {
		return b, err
	}
	b.InvoiceID = id
	return b, nil
}

func (a Adapter) writeBooking(b domain.Booking, w http.ResponseWriter) error {
	bs, err := json.Marshal(b)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bs)
	return nil
}

//=============================================================================
// Customer

func (a Adapter) readCustomer(r *http.Request, uid string) (domain.Customer, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return domain.Customer{}, err
	}
	var c domain.Customer
	if err := json.Unmarshal(body, &c); err != nil {
		return c, err
	}
	c.UserID = uid
	return c, nil
}

func (a Adapter) writeCustomer(c domain.Customer, w http.ResponseWriter) error {
	bs, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bs)
	return nil
}

//=============================================================================
// Invoice

func (a Adapter) readInvoice(r *http.Request) (domain.Invoice, error) {
	var i domain.Invoice
	// extract customerId from the URI
	cid, err := strconv.Atoi(mux.Vars(r)["customerId"])
	if err != nil {
		return i, err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return i, err
	}
	if err := json.Unmarshal(body, &i); err != nil {
		return i, err
	}
	i.CustomerID = cid
	return i, nil
}

func (a Adapter) writeInvoice(i domain.Invoice, w http.ResponseWriter) error {
	bs, err := json.Marshal(i)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bs)
	return nil
}

//=============================================================================
// Project

func (a Adapter) readProject(r *http.Request) (domain.Project, error) {
	var p domain.Project
	// extract customerId from the URI
	cid, err := strconv.Atoi(mux.Vars(r)["customerId"])
	if err != nil {
		return p, err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	if err := json.Unmarshal(body, &p); err != nil {
		return p, err
	}
	p.CustomerID = cid
	return p, nil
}

func (a Adapter) writeProject(p domain.Project, w http.ResponseWriter) error {
	bs, err := json.Marshal(p)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bs)
	return nil
}

//=============================================================================
// Rate

func (a Adapter) readRate(r *http.Request) (domain.Rate, error) {
	var rate domain.Rate
	// extract projectId from the URI
	pid, err := strconv.Atoi(mux.Vars(r)["projectId"])
	if err != nil {
		return rate, err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return rate, err
	}
	if err := json.Unmarshal(body, &rate); err != nil {
		return rate, err
	}
	rate.ProjectID = pid
	return rate, nil
}

//=============================================================================
// Handlers

// ActivitiesHandler returns a handler that knows how to retrieve activities
// for a user.
func (a Adapter) ActivitiesHandler(uc usecase.Activities) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		uid := a.currentUser(ctx)
		if len(uid) < 1 {
			w.WriteHeader(http.StatusUnauthorized)
		}
		// Runs the usecase to get the user's registered activities.
		as := uc.Run(uid)
		if len(as) < 1 {
			as = []domain.Activity{}
		}

		w.Header().Set("Content-Type", "application/json")

		// Mark response as chacheable by proxys and local caches.
		w.Header().Set("Cache-Control", "public, max-age=0")
		activities := NewActivitiesPresenter().Present(as)

		cc := r.Header.Get("Cache-Control")
		if len(cc) > 0 && strings.Contains(cc, "no-cache") {
			// Client requested a full refresh
			w.Header().Set("Last-Modified", activities.LastModified.Format(dateFormat))
			if _, err := w.Write(activities.Activities); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		// Respond with StatusNotModified if activities list has not been updated.
		if lms := r.Header.Get("Last-Modified-Since"); len(lms) > 0 {
			// Client sent conditional GET request, check mod date.
			mod, err := time.Parse(dateFormat, lms)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if truncateToSeconds(mod).Equal(truncateToSeconds(activities.LastModified)) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		// Respond activities list.
		w.Header().Set("Last-Modified", activities.LastModified.Format(dateFormat))
		if _, err := w.Write(activities.Activities); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// CreateActivityHandler returns a handler that knows how to create an activity.
func (a Adapter) CreateActivityHandler(uc usecase.CreateActivity) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		uid := a.currentUser(ctx)
		if len(uid) < 1 {
			w.WriteHeader(http.StatusUnauthorized)
		}
		act, err := a.readActivity(r, uid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Run the usecase to create an activity.
		created, err := uc.Run(act)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		location := fmt.Sprintf("%s/%d", r.URL.String(), created.ID)
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusCreated)
	}
}

// CreateBookingHandler returns a handler that knows how to create a booking.
func (a Adapter) CreateBookingHandler(uc usecase.CreateBooking) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		b, err := a.readBooking(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Runs the usecase to create a booking.
		created, err := uc.Run(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		location := fmt.Sprintf("%s/bookings/%d", r.URL.String(), created.ID)
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusCreated)
	}
}

// DeleteBookingHandler returns a handler that knows how to delete a booking.
func (a Adapter) DeleteBookingHandler(uc usecase.DeleteBooking) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["invoiceId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		bid, err := strconv.Atoi(mux.Vars(r)["bookingId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var b domain.Booking
		b.ID = bid
		b.InvoiceID = id
		// Runs the usecase to delete a booking.
		err = uc.Run(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// CreateCustomerHandler returns a handler that knows how to create a customer.
func (a Adapter) CreateCustomerHandler(uc usecase.CreateCustomer) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		uid := a.currentUser(ctx)
		if len(uid) < 1 {
			w.WriteHeader(http.StatusUnauthorized)
		}
		customer, err := a.readCustomer(r, uid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Runs the usecase to create a customer.
		created, err := uc.Run(customer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		location := fmt.Sprintf("%s/%d", r.URL.String(), created.ID)
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusCreated)
	}
}

// CreateInvoiceHandler returns a handler that knows how to create an invoice.
func (a Adapter) CreateInvoiceHandler(uc usecase.CreateInvoice) Handler {
	// A Closure that is closing over the createInvoice usecase instance.
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		i, err := a.readInvoice(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Runs the usecase to create an invoice.
		created, err := uc.Run(i)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		location := fmt.Sprintf("%s/%d", r.URL.String(), created.ID)
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusCreated)
	}
}

// GetInvoiceHandler returns a handler that knows how to deliver an invoice in
// either JSON or PDF format.
func (a Adapter) GetInvoiceHandler(uc usecase.GetInvoice) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		// extracts invoiceId from the URI
		id, err := strconv.Atoi(mux.Vars(r)["invoiceId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expand := ""
		q := r.URL.Query()
		// List of subressourcen that should be embedded in the invoice.
		if v, ok := q["expand"]; ok {
			// Use first element of the list.
			expand = v[0]
		}
		// JSON or PDF representation of the invoice.
		if p, ok := a.InvoicePresenter(w, r); ok {
			// Runs the usecase to get an invoice that optionaly includes and
			// lists subresources.
			i := uc.Run(id, expand)
			p.Present(NewHALInvoice(i))
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
	}
}

// UpdateInvoiceHandler returns a handler that knows how to update an ivoice.
func (a Adapter) UpdateInvoiceHandler(updateInvoice usecase.UpdateInvoice) Handler {
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		uid := a.currentUser(ctx)
		if len(uid) < 1 {
			w.WriteHeader(http.StatusUnauthorized)
		}
		// extract invoiceId from the URI
		id, err := strconv.Atoi(mux.Vars(r)["invoiceId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		i, err := a.readInvoice(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		i.ID = id
		// Runs the usecase to update an invoice.
		updateInvoice.Run(uid, i)
		w.WriteHeader(http.StatusNoContent)
	}
	return handler
}

// CreateProjectHandler returns a handler that knows how to create a project.
func (a Adapter) CreateProjectHandler(uc usecase.CreateProject) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		p, err := a.readProject(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Runs the usecase to create a project.
		created, err := uc.Run(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		location := fmt.Sprintf("%s/%d", r.URL.String(), created.ID)
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusCreated)
	}
}

// CreateRateHandler returns a handler that knows how to create a rate.
func (a Adapter) CreateRateHandler(uc usecase.CreateRate) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		rate, err := a.readRate(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Runs the usecase to create a rate.
		created, err := uc.Run(rate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		location := fmt.Sprintf("%s/activity/%d", r.URL.String(), created.ActivityID)
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusCreated)
	}
}
