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
	"github.com/tullo/invoice-mvp/usecase"
)

const (
	// admin user
	user       = "f8c39a31-9ced-4761-8a33-b9c628a67510"
	adminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8"
	customer   = 1
	pro1       = 1
	pro2       = 2
	inv1       = 1
)
const (
	act1 = iota + 1
	act2
	act3
)

func setupBaseData(r *database.FakeRepository) {
	// Customers
	r.CreateCustomer(domain.Customer{ID: customer, Name: "3skills", UserID: user})
	// Projects
	r.CreateProject(domain.Project{ID: pro1, Name: "Instanfoo.com", CustomerID: customer})
	r.CreateProject(domain.Project{ID: pro2, Name: "Covid19tracker.biz", CustomerID: customer})
	// Activities
	r.CreateActivity(domain.Activity{ID: act1, Name: "Programming", UserID: user})
	r.CreateActivity(domain.Activity{ID: act2, Name: "Quality control", UserID: user})
	r.CreateActivity(domain.Activity{ID: act3, Name: "Project management", UserID: user})
	// Project 1
	r.CreateRate(domain.Rate{ProjectID: pro1, ActivityID: act1, Price: 60}) // Programming
	r.CreateRate(domain.Rate{ProjectID: pro1, ActivityID: act2, Price: 55}) // Quality control
	// Project 2
	r.CreateRate(domain.Rate{ProjectID: pro2, ActivityID: act2, Price: 55}) // Quality control
	r.CreateRate(domain.Rate{ProjectID: pro2, ActivityID: act3, Price: 50}) // Project management
}

func booking(id, pid, aid int, h float32, d string) domain.Booking {
	return domain.Booking{
		InvoiceID:   id,
		ProjectID:   pid,
		ActivityID:  aid,
		Hours:       h,
		Description: d,
	}
}

func TestShouldUpdateState(t *testing.T) {
	// Setup
	r := database.NewFakeRepository()
	uc := usecase.NewUpdateInvoice(r)

	// Create invoice in "open" state
	i, err := r.CreateInvoice(domain.Invoice{ID: inv1, CustomerID: customer})
	if err != nil {
		t.Error(err)
	}
	// Update invoice state
	i.Status = "ready for aggregation"
	r.UpdateInvoice(i)

	// Run UpdateInvoice use case
	err = uc.Run(user, i)
	if err != nil {
		t.Error(err)
	}

	// Assert
	actual := r.GetInvoice(inv1)
	assert.Equal(t, "payment expected", actual.Status)
}

func TestAggregateBookings(t *testing.T) {
	//=========================================================================
	// Setup
	r := database.NewFakeRepository()
	setupBaseData(r)
	uc := usecase.NewUpdateInvoice(r)

	// Create bookings for project 1
	r.CreateBooking(booking(inv1, pro1, act1, 20, "Feature 4321 development"))
	r.CreateBooking(booking(inv1, pro1, act1, 12, "Rating impl"))
	r.CreateBooking(booking(inv1, pro1, act2, 3, "Rating test"))

	// Create bookings for project 2
	r.CreateBooking(booking(inv1, pro2, act3, 4, "Retrospective planing"))
	r.CreateBooking(booking(inv1, pro2, act3, 3, "Management offsite"))
	r.CreateBooking(booking(inv1, pro2, act2, 8, "Search testing"))

	// Create invoice in "open" state
	i, err := r.CreateInvoice(domain.Invoice{ID: inv1, CustomerID: customer})
	if err != nil {
		t.Error(err)
	}
	// advance invoice state
	i.Status = "ready for aggregation"
	r.UpdateInvoice(i)

	//=========================================================================
	// Run UpdateInvoice use case
	err = uc.Run(user, i)
	if err != nil {
		t.Error(err)
	}

	//=========================================================================
	// Assert
	mod, _ := time.Parse(time.RFC3339, "2020-11-20T12:00:00")
	status := "payment expected"
	expected := domain.Invoice{ID: 1, Status: status, CustomerID: customer}
	expected.AddPosition(pro1, "Programming", 32, 60)
	expected.AddPosition(pro1, "Quality control", 3, 55)
	expected.AddPosition(pro2, "Project management", 7, 50)
	expected.AddPosition(pro2, "Quality control", 8, 55)
	expected.Updated = mod

	actual := r.GetInvoice(inv1)
	actual.Updated = mod
	assert.Equal(t, expected, actual)
}

func TestHttpInvoiceAggregation(t *testing.T) {
	//=========================================================================
	// Setup
	r := database.NewFakeRepository()
	setupBaseData(r)
	updateInvoice := usecase.NewUpdateInvoice(r)

	// Create bookings for project 1
	r.CreateBooking(booking(inv1, pro1, act1, 20, "Feature 4321 development"))
	r.CreateBooking(booking(inv1, pro1, act1, 12, "Rating impl"))
	r.CreateBooking(booking(inv1, pro1, act2, 3, "Rating test"))

	// Create bookings for project 2
	r.CreateBooking(booking(inv1, pro2, act3, 4, "Retrospective planing"))
	r.CreateBooking(booking(inv1, pro2, act3, 3, "Management offsite"))
	r.CreateBooking(booking(inv1, pro2, act2, 8, "Search testing"))

	// Create invoice in "open" state
	i, err := r.CreateInvoice(domain.Invoice{ID: inv1, CustomerID: customer})
	if err != nil {
		t.Error(err)
	}
	// advance invoice state
	i.Status = "ready for aggregation"

	// Prepare HTTP-Request
	bs, _ := json.Marshal(&i)
	req, _ := http.NewRequest("PUT", "/customers/1/invoices/1", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

	//=========================================================================
	// Update invoice using PUT request
	res := httptest.NewRecorder()
	a := rest.NewAdapter()
	ui := a.UpdateInvoiceHandler(updateInvoice)
	ui = rest.JWTAuth(ui)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", ui).Methods("PUT")
	a.R.ServeHTTP(res, req)

	//=========================================================================
	// Assert
	mod, _ := time.Parse(time.RFC3339, "2020-11-20T12:00:00")
	status := "payment expected"
	expected := domain.Invoice{ID: 1, Status: status, CustomerID: customer}
	expected.AddPosition(pro1, "Programming", 32, 60)
	expected.AddPosition(pro1, "Quality control", 3, 55)
	expected.AddPosition(pro2, "Project management", 7, 50)
	expected.AddPosition(pro2, "Quality control", 8, 55)

	actual := r.GetInvoice(inv1)
	actual.Updated = mod
	assert.Equal(t, expected, actual)
}
