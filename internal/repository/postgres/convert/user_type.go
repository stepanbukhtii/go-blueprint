package convert

import (
	"github.com/samber/lo"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/model"
)

var UserType userType

type userType struct{}

func (userType) Domain(m *model.UserType) domain.UserType {
	return domain.UserType{
		Code:    m.Code,
		Name:    m.Name,
		IsAdmin: m.IsAdmin,
	}
}

func (c userType) DomainSlice(m model.UserTypeSlice) []domain.UserType {
	return lo.Map(m, func(m *model.UserType, _ int) domain.UserType { return c.Domain(m) })
}

func (userType) Model(ut *domain.UserType) *model.UserType {
	return &model.UserType{
		Code:    ut.Code,
		Name:    ut.Name,
		IsAdmin: ut.IsAdmin,
	}
}
