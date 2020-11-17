package domain

import "time"

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
