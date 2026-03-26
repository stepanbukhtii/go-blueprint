package request

import (
	"github.com/samber/lo"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
)

type CompanyURI struct {
	CompanyID string `uri:"company_id" binding:"required,uuid"`
}

type CreateCompany struct {
	Name      string `json:"name" binding:"required"`
	OwnerID   string `json:"owner_id" binding:"required"`
	ManagerID string `json:"manager_id"`
	LogoURL   string `json:"logo_url"`
}

func (r CreateCompany) ToDomain() domain.CreateCompanyInput {
	return domain.CreateCompanyInput{
		Name:      r.Name,
		OwnerID:   r.OwnerID,
		ManagerID: r.ManagerID,
		LogoURL:   r.LogoURL,
	}
}

type CreateCompanyMultiple []CreateCompany

func (r CreateCompanyMultiple) ToDomain() []domain.CreateCompanyInput {
	return lo.Map(r, func(r CreateCompany, _ int) domain.CreateCompanyInput { return r.ToDomain() })
}

type UpdateCompany struct {
	CompanyURI
	Name    string `json:"name"`
	LogoURL string `json:"logo_url"`
}

func (r UpdateCompany) ToDomain() domain.UpdateCompanyInput {
	return domain.UpdateCompanyInput{
		CompanyID: r.CompanyID,
		Name:      r.Name,
		LogoURL:   r.LogoURL,
	}
}
