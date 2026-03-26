package domain

import (
	"context"
	"errors"
	"time"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/shopspring/decimal"
	"github.com/stepanbukhtii/easy-tools/rest/api"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID               string
	Name             string
	Username         string
	Password         string
	PublicName       string
	Description      string
	UserType         UserType
	Age              int
	InitialAge       int
	Rate             float64
	LastRate         float64
	IsActive         bool
	ReadMessage      *bool
	Balance          decimal.Decimal
	LockBalance      *decimal.Decimal
	LastLogin        *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ManagerCompanies []Company
}

type UserRepository interface {
	Add(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Save(ctx context.Context, user *User) error
	Find(ctx context.Context, mods ...qm.QueryMod) (User, error)
	FindAll(ctx context.Context, mods ...qm.QueryMod) ([]User, error)
	FindAllPaginate(ctx context.Context, pagination api.Pagination, mods ...qm.QueryMod) ([]User, int64, error)
	Exists(ctx context.Context, mods ...qm.QueryMod) (bool, error)
	Remove(ctx context.Context, id string) error
}

type CreateUserInput struct {
	Name     string
	Username string
	Password string
}

type UpdateUserInput struct {
	UserID   string
	Name     string
	Username string
}

type UserService interface {
	ListPaginate(ctx context.Context, request api.Pagination) ([]User, int64, error)
	Get(ctx context.Context, id string) (User, error)
	Create(ctx context.Context, request CreateUserInput) (User, error)
	Update(ctx context.Context, request UpdateUserInput) (User, error)
	GeneratePublicName(ctx context.Context, id string) error
	UpdateRate(ctx context.Context, id string, lastRate float64) error
	Delete(ctx context.Context, id string) error
}

type UserAggregator interface {
	ListPaginate(ctx context.Context, request api.Pagination) ([]User, int64, error)
	Get(ctx context.Context, id string) (User, error)
}
