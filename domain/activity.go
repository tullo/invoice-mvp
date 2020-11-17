package domain

import (
	"fmt"
	"time"
)

// Activity represents a unit of work done by a user.
type Activity struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	UserID  string    `json:"userId"` // belongs to user
	Updated time.Time `json:"-"`
}

func (a Activity) String() string {
	return fmt.Sprintf("Id: %d Name: %s", a.ID, a.Name)
}
