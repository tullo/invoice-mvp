package usecase

import "github.com/tullo/invoice-mvp/domain"

// CreateRatePort is a small and use case specific interface.
type CreateRatePort interface {
	CreateRate(r domain.Rate) (domain.Rate, error)
}

// CreateRate implements the business logic.
type CreateRate struct {
	port CreateRatePort
}

// NewCreateRate instatiates the use case <Create Rate>.
func NewCreateRate(p CreateRatePort) CreateRate {
	return CreateRate{port: p}
}

// Run implements the use case <Create Rate>'.
func (u CreateRate) Run(r domain.Rate) (domain.Rate, error) {
	return u.port.CreateRate(r)
}
