package convert

import (
	"github.com/aarondl/opt/omit"
	"github.com/samber/lo"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/models"
)

var UserType userType

type userType struct{}

func (userType) Domain(m *models.UserType) domain.UserType {
	return domain.UserType{
		Code:    m.Code,
		Name:    m.Name,
		IsAdmin: m.IsAdmin,
	}
}

func (c userType) DomainSlice(m models.UserTypeSlice) []domain.UserType {
	return lo.Map(m, func(m *models.UserType, _ int) domain.UserType { return c.Domain(m) })
}

func (userType) Setter(ut *domain.UserType) *models.UserTypeSetter {
	return &models.UserTypeSetter{
		Code:    omit.From(ut.Code),
		Name:    omit.From(ut.Name),
		IsAdmin: omit.From(ut.IsAdmin),
	}
}
