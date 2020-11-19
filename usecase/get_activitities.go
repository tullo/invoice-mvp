package usecase

import "github.com/tullo/invoice-mvp/domain"

// ActivitiesPort is a small and use case specific interface.
type ActivitiesPort interface {
	Activities(userID string) []domain.Activity
}

// Activities implements the business logic.
type Activities struct {
	port ActivitiesPort
}

// NewActivities instatiates the use case <Get Activities>'.
func NewActivities(port ActivitiesPort) Activities {
	return Activities{port: port}
}

// Run implements the use case <Get Activities>'.
func (u Activities) Run(userID string) []domain.Activity {
	return u.port.Activities(userID)
}
