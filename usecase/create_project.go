package usecase

import "github.com/tullo/invoice-mvp/domain"

// CreateProjectPort is a small and use case specific interface.
type CreateProjectPort interface {
	CreateProject(p domain.Project) (domain.Project, error)
}

// CreateProject implements the business logic.
type CreateProject struct {
	port CreateProjectPort
}

// NewCreateProject instatiates the use case <Create Project>'.
func NewCreateProject(port CreateProjectPort) CreateProject {
	return CreateProject{port: port}
}

// Run implements the use case <Create Project>'.
func (u CreateProject) Run(project domain.Project) (domain.Project, error) {
	return u.port.CreateProject(project)
}
