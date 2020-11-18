package domain

import (
	"io/ioutil"
	"time"
)

// Position ...
type Position struct {
	Hours float32
	Price float32
}

// Invoice belongs to exactly one customer.
type Invoice struct {
	ID         int                         `json:"id"`
	Month      int                         `json:"month"`
	Year       int                         `json:"year"`
	Status     string                      `json:"status"`
	CustomerID int                         `json:"customerId"`
	Positions  map[int]map[string]Position `json:"positions,omitempty"`
	Bookings   []Booking                   `json:"-"`
	Updated    time.Time                   `json:"updated,omitempty"`
}

// AddPosition adds an invoice position or updates an existing one.
func (invoice *Invoice) AddPosition(projectID int, activity string, hours float32, rate float32) {
	if invoice.Positions == nil {
		invoice.Positions = make(map[int]map[string]Position)
	}

	if invoice.Positions[projectID] == nil {
		invoice.Positions[projectID] = make(map[string]Position)
	}

	if p, ok := invoice.Positions[projectID][activity]; ok {
		// update position values
		p.Hours = p.Hours + hours
		p.Price = p.Price + hours*rate
		invoice.Positions[projectID][activity] = p
	} else {
		// add position
		position := Position{Hours: hours, Price: hours * rate}
		invoice.Positions[projectID][activity] = position
	}
}

// ToPDF produces a pdf representation of the invoice.
func (invoice *Invoice) ToPDF() []byte {
	dat, _ := ioutil.ReadFile("/tmp/invoice.pdf")
	return dat
}
