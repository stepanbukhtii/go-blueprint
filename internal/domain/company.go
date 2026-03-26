package domain

import (
	"context"
	"errors"
	"time"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/stepanbukhtii/easy-tools/rest/api"
)

var (
	ErrCompanyNotFound = errors.New("user not found")
)

type Company struct {
	ID        string
	Name      string
	OwnerID   string
	ManagerID string
	IsActive  bool
	LogoURL   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Owner     *User
	Manager   *User
}

type CompanyRepository interface {
	Add(ctx context.Context, company *Company) error
	Update(ctx context.Context, company *Company) error
	Save(ctx context.Context, company *Company) error
	Find(ctx context.Context, mods ...qm.QueryMod) (Company, error)
	FindAll(ctx context.Context, mods ...qm.QueryMod) ([]Company, error)
	FindAllPaginate(ctx context.Context, pagination api.Pagination, mods ...qm.QueryMod) ([]Company, int64, error)
	Exists(ctx context.Context, mods ...qm.QueryMod) (bool, error)
	Remove(ctx context.Context, id string) error
}

type CreateCompanyInput struct {
	Name      string
	OwnerID   string
	ManagerID string
	LogoURL   string
}

type UpdateCompanyInput struct {
	CompanyID string
	Name      string
	LogoURL   string
}

type CompanyService interface {
	ListPaginate(ctx context.Context, request api.Pagination) ([]Company, int64, error)
	Get(ctx context.Context, id string) (Company, error)
	GetCompanyByOwner(ctx context.Context) ([]Company, error)
	Create(ctx context.Context, request CreateCompanyInput) (Company, error)
	CreateMultiple(ctx context.Context, request []CreateCompanyInput) ([]Company, error)
	Update(ctx context.Context, request UpdateCompanyInput) (Company, error)
	Delete(ctx context.Context, id string) error
}
