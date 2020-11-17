package domain

// Customer represents an entity that is related to one or more projects.
// A customer is owned by a user.
type Customer struct {
	ID       int       `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	UserID   string    `json:"userId,omitempty"`   // belongs to user
	Projects []Project `json:"projects,omitempty"` // has many projects
}
