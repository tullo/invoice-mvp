package usecase

import "github.com/tullo/invoice-mvp/domain"

// DeleteBookingPort is a small and use case specific interface.
type DeleteBookingPort interface {
	DeleteBooking(booking domain.Booking) error
}

// DeleteBooking implements the business logic.
type DeleteBooking struct {
	port DeleteBookingPort
}

// NewDeleteBooking instatiates the use case <Delete Booking>.
func NewDeleteBooking(port DeleteBookingPort) DeleteBooking {
	return DeleteBooking{port: port}
}

// Run implements the use case <Delete Booking>'.
func (u DeleteBooking) Run(booking domain.Booking) error {
	return u.port.DeleteBooking(booking)
}
