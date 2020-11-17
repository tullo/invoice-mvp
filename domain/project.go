package domain

// Project is an entity that belongs to excactly one customer.
type Project struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customerId"` // belongs to customer
	Name       string `json:"name"`
}
