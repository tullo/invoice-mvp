package database

import (
	"time"

	"github.com/tullo/invoice-mvp/domain"
)

// Repository is an in-memory store.
type Repository struct {
	customers  map[int]domain.Customer
	projects   map[int]domain.Project
	invoices   map[int]domain.Invoice
	bookings   map[int]domain.Booking
	activities map[int]domain.Activity
	rates      map[int]map[int]domain.Rate
}

// NewRepository creates a new repository.
func NewRepository() *Repository {
	r := Repository{
		customers:  make(map[int]domain.Customer),
		projects:   make(map[int]domain.Project),
		invoices:   make(map[int]domain.Invoice),
		bookings:   make(map[int]domain.Booking),
		activities: make(map[int]domain.Activity),
		rates:      make(map[int]map[int]domain.Rate),
	}

	return &r
}

// Customers gets all customers.
func (r *Repository) Customers() []domain.Customer {
	var cs []domain.Customer
	for _, c := range r.customers {
		cs = append(cs, c)
	}
	return cs
}

// Projects gets projects related to a customer.
func (r *Repository) Projects(customerID int) []domain.Project {
	var ps []domain.Project
	for _, p := range r.projects {
		if p.CustomerID == customerID {
			ps = append(ps, p)
		}
	}
	return ps
}

// Activities gets all activities.
func (r *Repository) Activities() []domain.Activity {
	var as []domain.Activity
	for _, c := range r.activities {
		as = append(as, c)
	}
	return as
}

// AddActivity adds an activity.
func (r *Repository) AddActivity(name string) int {
	a := domain.Activity{ID: r.nextActivityID(), Name: name}
	r.activities[a.ID] = a
	return a.ID
}

// AddCustomer adds a customer.
func (r *Repository) AddCustomer(name string) int {
	c := domain.Customer{ID: r.nextCustomerID(), Name: name}
	r.customers[c.ID] = c
	return c.ID
}

// AddProject add a project.
func (r *Repository) AddProject(name string, customerID int) int {
	p := domain.Project{ID: r.nextProjectID(), Name: name, CustomerID: customerID}
	r.projects[p.ID] = p
	return p.ID
}

// AddRate adds a rate.
func (r *Repository) AddRate(rate domain.Rate) {
	if _, ok := r.rates[rate.ProjectID]; !ok {
		r.rates[rate.ProjectID] = make(map[int]domain.Rate)
	}
	r.rates[rate.ProjectID][rate.ActivityID] = rate
}

// CreateInvoice creates an invoice.
func (r *Repository) CreateInvoice(i domain.Invoice) (domain.Invoice, error) {
	i.ID = r.nextInvoiceID()
	i.Status = "open"
	i.Bookings = []domain.Booking{}
	i.Updated = time.Now().UTC()
	r.invoices[i.ID] = i
	return i, nil
}

// CreateBooking creates a booking.
func (r *Repository) CreateBooking(b domain.Booking) (domain.Booking, error) {
	b.ID = r.nextBookingID()
	r.bookings[b.ID] = b
	return b, nil
}

// DeleteBooking deletes a booking.
func (r *Repository) DeleteBooking(id int) {
	delete(r.bookings, id)
}

// BookingsByInvoiceID finds bookings by invoice ID.
func (r *Repository) BookingsByInvoiceID(invoiceID int) []domain.Booking {
	var bs []domain.Booking
	for _, b := range r.bookings {
		if b.InvoiceID == invoiceID {
			bs = append(bs, b)
		}
	}
	return bs
}

// Update updates the invoice.
func (r *Repository) Update(i domain.Invoice) {
	r.invoices[i.ID] = i
}

// FindByID finds an invoice by its ID.
func (r *Repository) FindByID(id int) (domain.Invoice, bool) {
	i, ok := r.invoices[id]
	return i, ok
}

// ActivityByID ...
func (r *Repository) ActivityByID(id int) domain.Activity {
	return r.activities[id]
}

// GetInvoice gets an invoice by its ID.
func (r *Repository) GetInvoice(id int, join ...string) domain.Invoice {
	return r.invoices[id]
}

// RateByProjectIDAndActivityID gets rates mapped to a project and activity.
func (r *Repository) RateByProjectIDAndActivityID(projectID int, activityID int) domain.Rate {

	if _, ok := r.rates[projectID]; !ok {
		// project not found
		return domain.Rate{}
	}

	if _, ok := r.rates[projectID][activityID]; !ok {
		// activity not found
		return domain.Rate{}
	}

	return r.rates[projectID][activityID]
}

func (r *Repository) nextInvoiceID() int {
	nextID := 1
	for _, v := range r.invoices {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *Repository) nextCustomerID() int {
	nextID := 1
	for _, v := range r.customers {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *Repository) nextProjectID() int {
	nextID := 1
	for _, v := range r.projects {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *Repository) nextBookingID() int {
	nextID := 1
	for _, v := range r.bookings {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *Repository) nextActivityID() int {
	nextID := 1
	for _, v := range r.activities {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}
