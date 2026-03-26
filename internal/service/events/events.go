package events

const (
	UserCreatedEvent    = "user.created"
	UserUpdatedEvent    = "user.updated"
	CompanyUpdatedEvent = "company.updated"
)

type EventUserCreatedData struct {
	UserID string `json:"user_id"`
}

type EventUserUpdatedData struct {
	UserID string  `json:"user_id"`
	Rate   float64 `json:"rate"`
}

type EventCompanyUpdatedData struct {
	CompanyID string `json:"company_id"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
	ManagerID string `json:"manager_id"`
	IsActive  bool   `json:"is_active"`
	LogoURL   string `json:"logo_url"`
}
