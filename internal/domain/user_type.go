package domain

import (
	"context"

	"github.com/stepanbukhtii/easy-tools/errx"
	"github.com/stepanbukhtii/easy-tools/rest/api"
)

const (
	UserTypeDefault = "DEFAULT"
	UserTypeAdmin   = "ADMIN"
)

var (
	ErrUserTypeNotFound = errx.Wrap(api.ErrNotFound, "user type not found")
)

type UserType struct {
	Code    string
	Name    string
	IsAdmin bool
}

type UserTypeRepository interface {
	Add(ctx context.Context, user *UserType) error
	Update(ctx context.Context, user *UserType) error
	Find(ctx context.Context, code string) (UserType, error)
	FindAll(ctx context.Context) ([]UserType, error)
	Remove(ctx context.Context, id string) error
}

type CreateUserTypeInput struct {
	Code    string
	Name    string
	IsAdmin bool
}

type UpdateUserTypeInput struct {
	Code    string
	Name    string
	IsAdmin *bool
}

type UserTypeService interface {
	List(ctx context.Context) ([]UserType, error)
	Get(ctx context.Context, id string) (UserType, error)
	Create(ctx context.Context, request CreateUserTypeInput) (UserType, error)
	Update(ctx context.Context, request UpdateUserTypeInput) (UserType, error)
	Delete(ctx context.Context, id string) error
}
