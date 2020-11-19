package usecase

import "github.com/tullo/invoice-mvp/domain"

// CreateActivityPort is a small and use case specific interface.
type CreateActivityPort interface {
	CreateActivity(activity domain.Activity) (domain.Activity, error)
}

// CreateActivity implements the business logic.
type CreateActivity struct {
	port CreateActivityPort
}

// NewCreateActivity instatiates the use case <Create Activity>.
func NewCreateActivity(port CreateActivityPort) CreateActivity {
	return CreateActivity{port: port}
}

// Run implements the use case <Create Activity>'.
func (u CreateActivity) Run(activity domain.Activity) (domain.Activity, error) {
	return u.port.CreateActivity(activity)
}
