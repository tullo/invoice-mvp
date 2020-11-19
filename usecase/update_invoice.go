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
	BookingsByInvoiceID(invoiceID int) []domain.Booking
	// Gets the hourly rate used for an activiti on a specific project.
	RateByProjectIDAndActivityID(projectID int, activityID int) domain.Rate
	// Updates the invoice.
	UpdateInvoice(invoice domain.Invoice) error
}

// UpdateInvoice implements the business logic.
type UpdateInvoice struct {
	port UpdateInvoicePort
}

// NewUpdateInvoice instatiates the use case <Update Invoice>'.
func NewUpdateInvoice(port UpdateInvoicePort) UpdateInvoice {
	return UpdateInvoice{port: port}
}

// Run implements the use case <Update Invoice>'.
func (u UpdateInvoice) Run(uid string, invoice domain.Invoice) error {
	if invoice.IsReadyForAggregation() {
		bookings := u.port.BookingsByInvoiceID(invoice.ID)
		// Converts bookings to invoice positions.
		for _, b := range bookings {
			// Activity for the booking
			activity := u.port.ActivityByID(uid, b.ActivityID)
			// Hourly rate for the activity
			rate := u.port.RateByProjectIDAndActivityID(b.ProjectID, b.ActivityID)
			// Add invoice position with aggregated sum for the activity.
			invoice.AddPosition(b.ProjectID, activity.Name, b.Hours, rate.Price)
		}
		invoice.Status = "payment expected"
	}

	return u.port.UpdateInvoice(invoice)
}
