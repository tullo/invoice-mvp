package domain

// Rate is the price charged per hour for a specific activity done on a project.
type Rate struct {
	ProjectID  int     `json:"projectId"`
	ActivityID int     `json:"activityId"`
	Price      float32 `json:"price"`
}
