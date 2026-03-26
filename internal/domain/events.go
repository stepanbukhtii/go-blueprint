package domain

const (
	UserCreatedEvent = "user.created"
	UserUpdatedEvent = "user.updated"
)

type EventUserCreatedData struct {
	UserID string `json:"user_id"`
}

type EventUserUpdatedData struct {
	UserID string  `json:"user_id"`
	Rate   float64 `json:"rate"`
}
