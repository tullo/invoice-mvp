// This use case delegates the handling of invoices to its UpdateInvoicePort,
// which is a small Interface, that abstracts away the access to the actual
// database in use.

package usecase

import "github.com/tullo/invoice-mvp/domain"

// UpdateInvoicePort is a small and use case specific interface.
type UpdateInvoicePort interface {
	// Gets the activity, e.g. 'Programming'.
	ActivityByID(uid string, id int) domain.Activity
	// Gets the bookings on this invoice.
	BookingsByInvoiceID(id int) []domain.Booking
	// Gets the hourly rate used for an activiti on a specific project.
	RateByProjectIDAndActivityID(pid int, aid int) domain.Rate
	// Updates the invoice.
	UpdateInvoice(i domain.Invoice) error
}

// UpdateInvoice implements the business logic.
type UpdateInvoice struct {
	port UpdateInvoicePort
}

// NewUpdateInvoice instatiates the use case <Update Invoice>'.
func NewUpdateInvoice(p UpdateInvoicePort) UpdateInvoice {
	return UpdateInvoice{port: p}
}

// Run implements the use case <Update Invoice>'.
func (u UpdateInvoice) Run(uid string, i domain.Invoice) error {
	if i.IsReadyForAggregation() {
		bs := u.port.BookingsByInvoiceID(i.ID)
		// Converts bookings to invoice positions.
		for _, b := range bs {
			// Hourly rate for an activity on a project.
			r := u.port.RateByProjectIDAndActivityID(b.ProjectID, b.ActivityID)
			// Activity booked
			a := u.port.ActivityByID(uid, b.ActivityID)
			i.AddPosition(b.ProjectID, a.Name, b.Hours, r.Price)
		}
		i.Status = "payment expected"
	}

	return u.port.UpdateInvoice(i)
}
