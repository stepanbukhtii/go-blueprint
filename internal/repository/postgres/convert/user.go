package convert

import (
	"github.com/aarondl/null/v8"
	"github.com/samber/lo"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/model"
)

var User user

type user struct{}

func (user) Domain(m *model.User) domain.User {
	var managerCompanies []domain.Company
	if m.R != nil && m.R.ManagerCompanies != nil {
		managerCompanies = Company.DomainSlice(m.R.ManagerCompanies)
	}

	return domain.User{
		ID:               m.ID,
		Name:             m.Name,
		Username:         m.Username,
		Password:         m.Password,
		PublicName:       m.PublicName.String,
		Description:      m.Description.String,
		UserType:         domain.UserType{Code: m.UserType},
		Age:              m.Age,
		InitialAge:       m.InitialAge.Int,
		Rate:             m.Rate,
		LastRate:         m.LastRate.Float64,
		IsActive:         m.IsActive,
		ReadMessage:      m.ReadMessage.Ptr(),
		Balance:          m.Balance,
		LockBalance:      NullDecimalPtr(m.LockBalance),
		LastLogin:        m.LastLogin.Ptr(),
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		ManagerCompanies: managerCompanies,
	}
}

func (c user) DomainSlice(m model.UserSlice) []domain.User {
	return lo.Map(m, func(m *model.User, _ int) domain.User { return c.Domain(m) })
}

func (c user) Model(user *domain.User) *model.User {
	return &model.User{
		ID:          user.ID,
		Name:        user.Name,
		Username:    user.Username,
		Password:    user.Password,
		PublicName:  null.NewString(user.PublicName, user.PublicName != ""),
		Description: null.NewString(user.Description, user.Description != ""),
		UserType:    user.UserType.Code,
		Age:         user.Age,
		InitialAge:  null.NewInt(user.InitialAge, user.InitialAge != 0),
		Rate:        user.Rate,
		LastRate:    null.NewFloat64(user.Rate, user.Rate != 0),
		IsActive:    user.IsActive,
		ReadMessage: null.BoolFromPtr(user.ReadMessage),
		Balance:     user.Balance,
		LockBalance: DecimalPtr(user.LockBalance),
		LastLogin:   null.TimeFromPtr(user.LastLogin),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
