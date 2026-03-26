package request

import "github.com/stepanbukhtii/go-blueprint/internal/domain"

type UserTypeCodeURI struct {
	UserTypeCode string `uri:"user_type_code" binding:"required" swaggerignore:"true"`
}

type CreateUserType struct {
	Code    string `json:"code"  binding:"required"`
	Name    string `json:"name"  binding:"required"`
	IsAdmin bool   `json:"is_admin"  binding:"required"`
}

func (r CreateUserType) ToDomain() domain.CreateUserTypeInput {
	return domain.CreateUserTypeInput{
		Code:    r.Code,
		Name:    r.Name,
		IsAdmin: r.IsAdmin,
	}
}

type UpdateUserType struct {
	UserTypeCodeURI
	Name    string `json:"name"  binding:"required"`
	IsAdmin *bool  `json:"is_admin"  binding:"required"`
}

func (r UpdateUserType) ToDomain() domain.UpdateUserTypeInput {
	return domain.UpdateUserTypeInput{
		Code:    r.UserTypeCode,
		Name:    r.Name,
		IsAdmin: r.IsAdmin,
	}
}
