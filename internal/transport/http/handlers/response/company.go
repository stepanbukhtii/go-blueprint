package response

import (
	"github.com/samber/lo"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
)

type Company struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
	ManagerID string `json:"manager_id"`
	IsActive  bool   `json:"is_active"`
	LogoURL   string `json:"logo_url"`
}

func NewCompany(c domain.Company) Company {
	return Company{
		ID:        c.ID,
		Name:      c.Name,
		OwnerID:   c.OwnerID,
		ManagerID: c.ManagerID,
		IsActive:  c.IsActive,
		LogoURL:   c.LogoURL,
	}
}

func NewCompanies(companies []domain.Company) []Company {
	return lo.Map(companies, func(c domain.Company, index int) Company { return NewCompany(c) })
}
