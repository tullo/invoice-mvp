package domain

import (
	"io/ioutil"
	"time"
)

// Operation defines an operation on the invoice.
type Operation string

// Position models an invoice position.
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
	Bookings   []Booking                   `json:"bookings,omitempty"`
	Updated    time.Time                   `json:"updated,omitempty"`
	//Bookings   []Booking                 `json:"-"` excluded in json representation
}

// AddPosition adds an invoice position or updates an existing one.
func (invoice *Invoice) AddPosition(projectID int, activity string, hours float32, rate float32) {
	if invoice.Positions == nil {
		invoice.Positions = make(map[int]map[string]Position)
	}

	// Instantiate positions map for a project.
	if invoice.Positions[projectID] == nil {
		invoice.Positions[projectID] = make(map[string]Position)
	}

	// Create or update a position for an activity on a project.
	if p, ok := invoice.Positions[projectID][activity]; ok {
		// update aggregated position sum values for the activity.
		p.Hours = p.Hours + hours
		p.Price = p.Price + hours*rate
		invoice.Positions[projectID][activity] = p
	} else {
		// add position for the activity.
		position := Position{Hours: hours, Price: hours * rate}
		invoice.Positions[projectID][activity] = position
	}
}

// IsReadyForAggregation indicates whether an invoice is in
// "ready for aggregation" state.
func (invoice Invoice) IsReadyForAggregation() bool {
	return invoice.Status == "ready for aggregation"
}

// ToPDF produces a pdf representation of the invoice.
func (invoice *Invoice) ToPDF() []byte {
	dat, _ := ioutil.ReadFile("/tmp/invoice.pdf")
	return dat
}

// Operations returns allowed operations depending on current invoice state.
func (invoice Invoice) Operations() []Operation {
	switch invoice.Status {
	case "open":
		return []Operation{"book", "charge", "cancel", "bookings"}
	case "payment expected":
		return []Operation{"payment", "bookings"}
	case "paid":
		return []Operation{"archive"}
	case "archived":
		return []Operation{"revoke"}
	case "revoked":
		return []Operation{"archive"}
	default:
		return []Operation{}
	}
}
