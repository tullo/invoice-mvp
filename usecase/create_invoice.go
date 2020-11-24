package usecase

import "github.com/tullo/invoice-mvp/domain"

// CreateInvoicePort is a small and use case specific interface.
type CreateInvoicePort interface {
	CreateInvoice(i domain.Invoice) (domain.Invoice, error)
}

// CreateInvoice implements the business logic.
type CreateInvoice struct {
	port CreateInvoicePort
}

// NewCreateInvoice instatiates the use case <Create Invoice>'.
func NewCreateInvoice(p CreateInvoicePort) CreateInvoice {
	return CreateInvoice{port: p}
}

// Run implements the use case <Create Invoice>'.
func (u CreateInvoice) Run(i domain.Invoice) (domain.Invoice, error) {
	return u.port.CreateInvoice(i)
}
