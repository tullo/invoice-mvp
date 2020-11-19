package usecase

import "github.com/tullo/invoice-mvp/domain"

// CreateInvoicePort is a small and use case specific interface.
type CreateInvoicePort interface {
	CreateInvoice(invoice domain.Invoice) (domain.Invoice, error)
}

// CreateInvoice implements the business logic.
type CreateInvoice struct {
	port CreateInvoicePort
}

// NewCreateInvoice instatiates the use case <Create Invoice>'.
func NewCreateInvoice(port CreateInvoicePort) CreateInvoice {
	return CreateInvoice{port: port}
}

// Run implements the use case <Create Invoice>'.
func (u CreateInvoice) Run(invoice domain.Invoice) (domain.Invoice, error) {
	return u.port.CreateInvoice(invoice)
}
