package usecase

import "github.com/tullo/invoice-mvp/domain"

// CreateCustomerPort is a small and use case specific interface.
type CreateCustomerPort interface {
	CreateCustomer(customer domain.Customer) (domain.Customer, error)
}

// CreateCustomer implements the business logic.
type CreateCustomer struct {
	port CreateCustomerPort
}

// NewCreateCustomer instatiates the use case <Create Customer>'.
func NewCreateCustomer(port CreateCustomerPort) CreateCustomer {
	return CreateCustomer{port: port}
}

// Run implements the use case <Create Customer>'.
func (u CreateCustomer) Run(Customer domain.Customer) (domain.Customer, error) {
	return u.port.CreateCustomer(Customer)
}
