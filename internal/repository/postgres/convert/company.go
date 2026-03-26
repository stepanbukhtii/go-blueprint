package convert

import (
	"github.com/aarondl/null/v8"
	"github.com/samber/lo"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/model"
)

var Company company

type company struct{}

func (company) Domain(m *model.Company) domain.Company {
	return domain.Company{
		ID:        m.ID,
		Name:      m.Name,
		OwnerID:   m.OwnerID,
		ManagerID: m.ManagerID.String,
		IsActive:  m.IsActive,
		LogoURL:   m.LogoURL.String,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (c company) DomainSlice(m model.CompanySlice) []domain.Company {
	return lo.Map(m, func(m *model.Company, _ int) domain.Company { return c.Domain(m) })
}

func (c company) Model(company *domain.Company) *model.Company {
	return &model.Company{
		ID:        company.ID,
		Name:      company.Name,
		OwnerID:   company.OwnerID,
		ManagerID: null.NewString(company.ManagerID, company.ManagerID != ""),
		IsActive:  company.IsActive,
		LogoURL:   null.NewString(company.LogoURL, company.LogoURL != ""),
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}
}
