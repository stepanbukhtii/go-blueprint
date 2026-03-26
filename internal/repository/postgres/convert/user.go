package convert

import (
	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/samber/lo"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/models"
)

var User user

type user struct{}

func (user) Domain(m *models.User) domain.User {
	return domain.User{
		ID:          m.ID,
		Name:        m.Name,
		Username:    m.Username,
		Password:    m.Password,
		PublicName:  m.PublicName.GetOrZero(),
		Description: m.Description.GetOrZero(),
		UserType:    domain.UserType{Code: m.UserType},
		Age:         int(m.Age),
		InitialAge:  int(m.InitialAge.GetOrZero()),
		Rate:        m.Rate,
		LastRate:    m.LastRate.GetOrZero(),
		IsActive:    m.IsActive,
		ReadMessage: m.ReadMessage.Ptr(),
		Balance:     m.Balance,
		LockBalance: m.LockBalance.Ptr(),
		LastLogin:   m.LastLogin.Ptr(),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (c user) DomainSlice(m models.UserSlice) []domain.User {
	return lo.Map(m, func(m *models.User, _ int) domain.User { return c.Domain(m) })
}

func (user) Setter(u *domain.User) *models.UserSetter {
	return &models.UserSetter{
		Name:        omit.From(u.Name),
		Username:    omit.From(u.Username),
		Password:    omit.From(u.Password),
		PublicName:  omitnull.FromNull(null.FromCond(u.PublicName, u.PublicName != "")),
		Description: omitnull.FromNull(null.FromCond(u.Description, u.Description != "")),
		UserType:    omit.From(u.UserType.Code),
		Age:         omit.From(int32(u.Age)),
		InitialAge:  omitnull.FromNull(null.FromCond(int32(u.InitialAge), u.InitialAge != 0)),
		Rate:        omit.From(u.Rate),
		LastRate:    omitnull.FromNull(null.FromCond(u.LastRate, u.LastRate != 0)),
		IsActive:    omit.From(u.IsActive),
		ReadMessage: omitnull.FromPtr(u.ReadMessage),
		Balance:     omit.From(u.Balance),
		LockBalance: omitnull.FromPtr(u.LockBalance),
		LastLogin:   omitnull.FromPtr(u.LastLogin),
	}
}
