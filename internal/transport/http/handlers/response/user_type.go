package response

import (
	"github.com/samber/lo"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
)

type UserType struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

func NewUserType(ut domain.UserType) UserType {
	return UserType{
		Code:    ut.Code,
		Name:    ut.Name,
		IsAdmin: ut.IsAdmin,
	}
}

func NewUserTypes(userTypes []domain.UserType) []UserType {
	return lo.Map(userTypes, func(ut domain.UserType, index int) UserType { return NewUserType(ut) })
}
