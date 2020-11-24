package usecase

import "github.com/tullo/invoice-mvp/domain"

// DeleteBookingPort is a small and use case specific interface.
type DeleteBookingPort interface {
	DeleteBooking(b domain.Booking) error
}

// DeleteBooking implements the business logic.
type DeleteBooking struct {
	port DeleteBookingPort
}

// NewDeleteBooking instatiates the use case <Delete Booking>.
func NewDeleteBooking(p DeleteBookingPort) DeleteBooking {
	return DeleteBooking{port: p}
}

// Run implements the use case <Delete Booking>'.
func (u DeleteBooking) Run(b domain.Booking) error {
	return u.port.DeleteBooking(b)
}
