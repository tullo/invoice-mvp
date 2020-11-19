package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/tullo/invoice-mvp/domain"
)

// FakeRepository is an in-memory store.
type FakeRepository struct {
	activities map[string]map[int]domain.Activity
	bookings   map[int]map[int]domain.Booking
	customers  map[int]domain.Customer
	invoices   map[int]domain.Invoice
	projects   map[int]domain.Project
	rates      map[int]map[int]domain.Rate
}

// NewFakeRepository creates a new repository.
func NewFakeRepository() *FakeRepository {
	r := FakeRepository{
		activities: make(map[string]map[int]domain.Activity),
		bookings:   make(map[int]map[int]domain.Booking),
		customers:  make(map[int]domain.Customer),
		invoices:   make(map[int]domain.Invoice),
		projects:   make(map[int]domain.Project),
		rates:      make(map[int]map[int]domain.Rate),
	}

	return &r
}

//=============================================================================
// Activities

// Activities gets all activities.
func (r *FakeRepository) Activities(userID string) []domain.Activity {
	var as []domain.Activity
	for _, a := range r.activities[userID] {
		as = append(as, a)
	}
	return as
}

// ActivityByID gets an users activity.
func (r *FakeRepository) ActivityByID(uid string, aid int) domain.Activity {
	return r.activities[uid][aid]
}

// CreateActivity adds an activity to a users activity list.
func (r *FakeRepository) CreateActivity(a domain.Activity) (domain.Activity, error) {
	a.ID = r.nextActivityID(a.UserID)
	a.Updated = time.Now().UTC()
	if _, ok := r.activities[a.UserID]; !ok {
		// initiate activities map
		am := make(map[int]domain.Activity)
		r.activities[a.UserID] = am
	}
	am := r.activities[a.UserID]
	am[a.ID] = a
	r.activities[a.UserID] = am
	return a, nil
}

//=============================================================================
// Bookings

// BookingsByInvoiceID finds bookings by invoice ID.
func (r *FakeRepository) BookingsByInvoiceID(invoiceID int) []domain.Booking {
	var bs []domain.Booking
	if bm, ok := r.bookings[invoiceID]; ok {
		for _, b := range bm {
			bs = append(bs, b)
		}
	}
	return bs
}

// CreateBooking creates a booking.
func (r *FakeRepository) CreateBooking(b domain.Booking) (domain.Booking, error) {
	b.ID = r.nextBookingID(b.InvoiceID)
	if bs, ok := r.bookings[b.InvoiceID]; ok {
		bs[b.ID] = b
	} else {
		bm := make(map[int]domain.Booking)
		bm[b.ID] = b
		r.bookings[b.InvoiceID] = bm
	}
	return b, nil

}

// DeleteBooking deletes a booking.
func (r *FakeRepository) DeleteBooking(b domain.Booking) error {
	if bm, ok := r.bookings[b.InvoiceID]; ok {
		if _, ok := bm[b.ID]; ok {
			delete(bm, b.ID)
			return nil
		}
	}
	return fmt.Errorf("failed to delete booking %d on invoice %d", b.ID, b.InvoiceID)
}

//=============================================================================
// Customers

// CreateCustomer adds a customer to the repository.
func (r *FakeRepository) CreateCustomer(c domain.Customer) (domain.Customer, error) {
	c.ID = r.nextCustomerID()
	r.customers[c.ID] = c
	return c, nil

}

// Customers gets all customers.
func (r *FakeRepository) Customers() []domain.Customer {
	var cs []domain.Customer
	for _, c := range r.customers {
		cs = append(cs, c)
	}
	return cs
}

// CustomerByID finds a customer by customer ID.
func (r *FakeRepository) CustomerByID(id int) domain.Customer {
	return r.customers[id]
}

//=============================================================================
// Invoices

// GetInvoice gets an invoice by its ID and optionally embeds bookings.
func (r *FakeRepository) GetInvoice(id int, join ...string) domain.Invoice {
	i := r.invoices[id]
	if len(join) > 0 {
		if strings.Contains(join[0], "bookings") {
			i.Bookings = r.BookingsByInvoiceID(id)
		}
	}
	return i
}

// CreateInvoice creates an invoice in the repository.
func (r *FakeRepository) CreateInvoice(i domain.Invoice) (domain.Invoice, error) {
	var bs []domain.Booking
	i.ID = r.nextInvoiceID()
	i.Status = "open"
	i.Bookings = bs
	i.Updated = time.Now().UTC()
	r.invoices[i.ID] = i
	return i, nil
}

// UpdateInvoice updates the invoice in the repository.
func (r *FakeRepository) UpdateInvoice(i domain.Invoice) error {
	r.invoices[i.ID] = i
	return nil
}

//=============================================================================
// Projects

// CreateProject creates a project in the repository.
func (r *FakeRepository) CreateProject(p domain.Project) (domain.Project, error) {
	p.ID = r.nextProjectID()
	r.projects[p.ID] = p
	return p, nil
}

// ProjectByID finds a project by project ID.
func (r *FakeRepository) ProjectByID(id int) domain.Project {
	return r.projects[id]
}

// Projects gets projects related to a customer.
func (r *FakeRepository) Projects(customerID int) []domain.Project {
	var ps []domain.Project
	for _, p := range r.projects {
		if p.CustomerID == customerID {
			ps = append(ps, p)
		}
	}
	return ps
}

//=============================================================================
// Rates

// CreateRate creates a rate in the repository.
func (r *FakeRepository) CreateRate(rate domain.Rate) (domain.Rate, error) {
	if _, ok := r.rates[rate.ProjectID]; !ok {
		// create map for project rates
		rates := make(map[int]domain.Rate)
		r.rates[rate.ProjectID] = rates
	}
	r.rates[rate.ProjectID][rate.ActivityID] = rate
	return rate, nil
}

// RateByProjectIDAndActivityID gets the rate mapped to a project and ID.
func (r *FakeRepository) RateByProjectIDAndActivityID(projectID int, activityID int) domain.Rate {
	var rate domain.Rate
	if _, ok := r.rates[projectID]; !ok {
		// project not found
		return rate
	}

	if _, ok := r.rates[projectID][activityID]; !ok {
		// activity not found
		return rate
	}

	return r.rates[projectID][activityID]
}

func (r *FakeRepository) nextInvoiceID() int {
	nextID := 1
	for _, v := range r.invoices {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *FakeRepository) nextCustomerID() int {
	nextID := 1
	for _, v := range r.customers {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *FakeRepository) nextProjectID() int {
	nextID := 1
	for _, v := range r.projects {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *FakeRepository) nextBookingID(inv int) int {
	nextID := 1
	if _, ok := r.bookings[inv]; ok {
		for _, v := range r.bookings[inv] {
			if v.ID >= nextID {
				nextID = v.ID + 1
			}
		}
	}
	return nextID
}

func (r *FakeRepository) nextActivityID(uid string) int {
	nextID := 1
	if _, ok := r.activities[uid]; ok {
		for _, v := range r.activities[uid] {
			if v.ID >= nextID {
				nextID = v.ID + 1
			}
		}
	}
	return nextID
}
