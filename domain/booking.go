package domain

import "fmt"

// Booking is a record for work done on excatly one project.
type Booking struct {
	ID          int     `json:"-"`
	Day         int     `json:"day"`
	Hours       float32 `json:"hours"`
	Description string  `json:"description"`
	InvoiceID   int     `json:"invoiceId"`            // belongs to invoice
	ProjectID   int     `json:"projectId,omitempty"`  // belongs to project
	ActivityID  int     `json:"activityId,omitempty"` // belongs to activity
}

func (b Booking) String() string {
	return fmt.Sprintf("Id: %d Day: %d Hours: %f Description: %s InvoiceID: %d ProjectID: %d ActivityID: %d",
		b.ID, b.Day, b.Hours, b.Description, b.InvoiceID, b.ProjectID, b.ActivityID)
}
