package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"

	"github.com/stepanbukhtii/easy-tools/errx"
	"github.com/stepanbukhtii/easy-tools/rest/api"
)

var (
	ErrUserNotFound     = errx.Wrap(api.ErrNotFound, "user not found")
	ErrUsernameIsExists = errx.Wrap(api.ErrConflict, "username is exists")
)

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
	Find(ctx context.Context, mods ...bob.Mod[*dialect.SelectQuery]) (User, error)
	FindAll(ctx context.Context, mods ...bob.Mod[*dialect.SelectQuery]) ([]User, error)
	FindAllPaginate(
		ctx context.Context,
		pagination api.Pagination,
		mods ...bob.Mod[*dialect.SelectQuery],
	) ([]User, int64, error)
	Exists(ctx context.Context, mods ...bob.Mod[*dialect.SelectQuery]) (bool, error)
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
