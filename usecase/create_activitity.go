package usecase

import "github.com/tullo/invoice-mvp/domain"

// CreateActivityPort is a small and use case specific interface.
type CreateActivityPort interface {
	CreateActivity(a domain.Activity) (domain.Activity, error)
}

// CreateActivity implements the business logic.
type CreateActivity struct {
	port CreateActivityPort
}

// NewCreateActivity instatiates the use case <Create Activity>.
func NewCreateActivity(p CreateActivityPort) CreateActivity {
	return CreateActivity{port: p}
}

// Run implements the use case <Create Activity>'.
func (u CreateActivity) Run(a domain.Activity) (domain.Activity, error) {
	return u.port.CreateActivity(a)
}
