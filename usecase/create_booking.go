// This use case delegates the handling of bookings to its CreateBookingPort,
// which is a small Interface, that abstracts away the access to the actual
// database in use.

package usecase

import "github.com/tullo/invoice-mvp/domain"

// CreateBookingPort is a small and use case specific interface.
type CreateBookingPort interface {
	CreateBooking(booking domain.Booking) (domain.Booking, error)
}

// CreateBooking implements the business logic.
type CreateBooking struct {
	port CreateBookingPort
}

// NewCreateBooking instatiates the use case <Create Booking>.
func NewCreateBooking(port CreateBookingPort) CreateBooking {
	return CreateBooking{port: port}
}

// Run implements the use case <Create Booking>'.
func (u CreateBooking) Run(booking domain.Booking) (domain.Booking, error) {
	return u.port.CreateBooking(booking)
}
