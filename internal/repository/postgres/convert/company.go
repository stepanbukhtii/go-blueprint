package convert

import (
	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/samber/lo"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/models"
)

var Company company

type company struct{}

func (company) Domain(m *models.Company) domain.Company {
	return domain.Company{
		ID:        m.ID,
		Name:      m.Name,
		OwnerID:   m.OwnerID,
		ManagerID: m.ManagerID.GetOrZero(),
		IsActive:  m.IsActive,
		LogoURL:   m.LogoURL.GetOrZero(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (c company) DomainSlice(m models.CompanySlice) []domain.Company {
	return lo.Map(m, func(m *models.Company, _ int) domain.Company { return c.Domain(m) })
}

func (company) Setter(company *domain.Company) *models.CompanySetter {
	return &models.CompanySetter{
		Name:      omit.From(company.Name),
		OwnerID:   omit.From(company.OwnerID),
		ManagerID: omitnull.FromNull(null.FromCond(company.ManagerID, company.ManagerID != "")),
		IsActive:  omit.From(company.IsActive),
		LogoURL:   omitnull.FromNull(null.FromCond(company.LogoURL, company.LogoURL != "")),
	}
}
