package usecase_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tullo/invoice-mvp/database"
	"github.com/tullo/invoice-mvp/domain"
	"github.com/tullo/invoice-mvp/rest"
	"github.com/tullo/invoice-mvp/roles"
	"github.com/tullo/invoice-mvp/usecase"
)

func TestHttpAddBooking(t *testing.T) {
	//=========================================================================
	// Setup
	r := database.NewFakeRepository()
	setupBaseData(r)
	createBooking := usecase.NewCreateBooking(r)

	// Create invoice in "open" state
	_, err := r.CreateInvoice(domain.Invoice{ID: inv1, CustomerID: customer})
	if err != nil {
		t.Error(err)
	}

	// Prepare HTTP-Request
	b := domain.Booking{
		Day:         31,
		Hours:       2.5,
		Description: "Front: bugfix #6789",
		ProjectID:   pro1,
		ActivityID:  act1,
	}
	bs, _ := json.Marshal(&b)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/book/%v", inv1), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

	//=========================================================================
	// Add booking using POST request
	res := httptest.NewRecorder()
	a := rest.NewAdapter()
	cb := a.CreateBookingHandler(createBooking)
	cb = rest.JWTAuth(roles.AssertOwnsInvoice(cb, r))
	a.Handle("/book/{invoiceId:[0-9]+}", cb).Methods("POST")
	a.R.ServeHTTP(res, req)

	//=========================================================================
	// Assert
	assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
	assert.Equal(t, "/book/1/bookings/1", res.Result().Header["Location"][0])
	mod, _ := time.Parse(time.RFC3339, "2020-11-20T12:00:00")
	status := "open"
	expected := domain.Invoice{ID: inv1, Status: status, CustomerID: customer}
	expected.Bookings = append(expected.Bookings, domain.Booking{
		ID:          1,
		Day:         b.Day,
		Hours:       b.Hours,
		Description: b.Description,
		ProjectID:   b.ProjectID,
		ActivityID:  b.ActivityID,
		InvoiceID:   inv1,
	})
	expansion := "bookings"
	actual := r.GetInvoice(inv1, expansion)
	actual.Updated = mod
	assert.Equal(t, expected, actual)
}
