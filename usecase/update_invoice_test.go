package usecase_test

import (
	"bytes"
	"encoding/json"
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
	user     = "1234"
	customer = 1
	pro1     = 1
	pro2     = 2
)
const (
	act1 = iota + 1
	act2
	act3
)

func setupBaseData(r *database.FakeRepository) {
	// Projects
	r.CreateProject(domain.Project{ID: pro1, Name: "Instanfoo.com", CustomerID: customer})
	r.CreateProject(domain.Project{ID: pro2, Name: "Covid19tracker.biz", CustomerID: customer})
	// Activities
	r.CreateActivity(domain.Activity{ID: 1, Name: "Programming", UserID: user})
	r.CreateActivity(domain.Activity{ID: 2, Name: "Quality control", UserID: user})
	r.CreateActivity(domain.Activity{ID: 3, Name: "Project management", UserID: user})
	// Project 1
	r.CreateRate(domain.Rate{ProjectID: pro1, ActivityID: act1, Price: 60}) // Programming
	r.CreateRate(domain.Rate{ProjectID: pro1, ActivityID: act2, Price: 55}) // Quality control
	// Project 2
	r.CreateRate(domain.Rate{ProjectID: pro2, ActivityID: act2, Price: 55}) // Quality control
	r.CreateRate(domain.Rate{ProjectID: pro2, ActivityID: act3, Price: 50}) // Project management
}

func TestShouldUpdateState(t *testing.T) {
	// Setup
	r := database.NewFakeRepository()
	uc := usecase.NewUpdateInvoice(r)

	// Create invoice in "open" state
	i, err := r.CreateInvoice(domain.Invoice{ID: 1, CustomerID: customer})
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
	actual := r.GetInvoice(1)
	assert.Equal(t, "payment expected", actual.Status)
}

func TestAggregateBookings(t *testing.T) {
	// Setup
	r := database.NewFakeRepository()
	setupBaseData(r)
	uc := usecase.NewUpdateInvoice(r)

	// Create bookings for project 1
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro1,
		ActivityID: act1, Hours: 20, Description: "Feature 4321 development"})
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro1,
		ActivityID: act1, Hours: 12, Description: "Rating impl"})
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro1,
		ActivityID: act2, Hours: 3, Description: "Rating test"})

	// Create bookings for project 2
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro2,
		ActivityID: act3, Hours: 4, Description: "Retrospective planing"})
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro2,
		ActivityID: act3, Hours: 3, Description: "Management offsite"})
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro2,
		ActivityID: act2, Hours: 8, Description: "Search testing"})

	// Create invoice in "open" state
	i, err := r.CreateInvoice(domain.Invoice{ID: 1, CustomerID: customer})
	if err != nil {
		t.Error(err)
	}
	// Update invoice state
	mod, _ := time.Parse(time.RFC3339, "2020-11-20T12:00:00")
	i.Status = "ready for aggregation"
	i.Updated = mod
	r.UpdateInvoice(i)

	// Run UpdateInvoice use case
	err = uc.Run(user, i)
	if err != nil {
		t.Error(err)
	}

	// Assert
	status := "payment expected"
	expected := domain.Invoice{ID: 1, Status: status, CustomerID: customer}
	expected.AddPosition(pro1, "Programming", 32, 60)
	expected.AddPosition(pro1, "Quality control", 3, 55)
	expected.AddPosition(pro2, "Project management", 7, 50)
	expected.AddPosition(pro2, "Quality control", 8, 55)
	expected.Updated = mod

	actual := r.GetInvoice(1)
	assert.Equal(t, expected, actual)
}

func TestHttpInvoiceAggregation(t *testing.T) {
	// Setup
	r := database.NewFakeRepository()
	setupBaseData(r)
	uc := usecase.NewUpdateInvoice(r)

	// Create bookings for project 1
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro1,
		ActivityID: act1, Hours: 20, Description: "Feature 4321 development"})
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro1,
		ActivityID: act1, Hours: 12, Description: "Rating impl"})
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro1,
		ActivityID: act2, Hours: 3, Description: "Rating test"})

	// Create bookings for project 2
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro2,
		ActivityID: act3, Hours: 4, Description: "Retrospective planing"})
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro2,
		ActivityID: act3, Hours: 3, Description: "Management offsite"})
	r.CreateBooking(domain.Booking{InvoiceID: 1, ProjectID: pro2,
		ActivityID: act2, Hours: 8, Description: "Search testing"})

	i, err := r.CreateInvoice(domain.Invoice{ID: 1, CustomerID: customer})
	if err != nil {
		t.Error(err)
	}
	// Update invoice state
	i.Status = "ready for aggregation"

	// Prepare HTTP-Request
	bs, _ := json.Marshal(&i)
	req, _ := http.NewRequest("PUT", "/customers/1/invoices/1", bytes.NewReader(bs))

	// Run
	res := httptest.NewRecorder()
	a := rest.NewAdapter()
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", a.UpdateInvoiceHandler(uc)).Methods("PUT")
	a.R.ServeHTTP(res, req)

	// Assert
	status := "payment expected"
	expected := domain.Invoice{ID: 1, Status: status, CustomerID: customer}
	expected.AddPosition(pro1, "Programming", 32, 60)
	expected.AddPosition(pro1, "Quality control", 3, 55)
	expected.AddPosition(pro2, "Project management", 7, 50)
	expected.AddPosition(pro2, "Quality control", 8, 55)

	actual := r.GetInvoice(1)
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.CustomerID, actual.CustomerID)
	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.Bookings, actual.Bookings)
	assert.Equal(t, expected.Positions, actual.Positions)
}
