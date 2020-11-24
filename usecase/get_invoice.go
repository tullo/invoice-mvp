package usecase

import "github.com/tullo/invoice-mvp/domain"

// GetInvoicePort is a small and use case specific interface.
type GetInvoicePort interface {
	GetInvoice(id int, join ...string) domain.Invoice
}

// GetInvoice implements the business logic.
type GetInvoice struct {
	port GetInvoicePort
}

// NewGetInvoice instatiates the use case <Get Invoice>'.
func NewGetInvoice(p GetInvoicePort) GetInvoice {
	return GetInvoice{port: p}
}

// Run implements the use case <Get Invoice>'.
func (u GetInvoice) Run(id int, join string) domain.Invoice {
	return u.port.GetInvoice(id, join)
}
