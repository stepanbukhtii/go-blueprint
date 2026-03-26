package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/redis/go-redis/v9"
	"github.com/stepanbukhtii/easy-tools/cache"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/cached"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres"
)

type Repository interface {
	Company() domain.CompanyRepository
	User() domain.UserRepository
	UserType() domain.UserTypeRepository
	RunInTransaction(ctx context.Context, exec func(tx Repository) error) error
}

var ErrAlreadyInTransaction = errors.New("already in transaction")

type repository struct {
	serviceName string
	db          *sql.DB
	tx          *sql.Tx
	redis       *redis.Client
}

func NewRepository(serviceName string, db *sql.DB, redis *redis.Client) Repository {
	return &repository{
		serviceName: serviceName,
		db:          db,
		redis:       redis,
	}
}

func (r *repository) newRepositoryTx(tx *sql.Tx) Repository {
	return &repository{
		serviceName: r.serviceName,
		db:          r.db,
		tx:          tx,
		redis:       r.redis,
	}
}

func (r *repository) User() domain.UserRepository {
	return postgres.NewUser(r.exec())
}

func (r *repository) Company() domain.CompanyRepository {
	return postgres.NewCompany(r.exec())
}

func (r *repository) UserType() domain.UserTypeRepository {
	c := cache.NewRedisMap[domain.UserType](r.redis, r.serviceName, "user_type", cache.DefaultTTL)
	return cached.NewUserTypeRepository(c, postgres.NewUserType(r.exec()))
}

func (r *repository) RunInTransaction(ctx context.Context, exec func(tx Repository) error) error {
	if r.tx != nil {
		return ErrAlreadyInTransaction
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		recoveredFrom := recover()
		if recoveredFrom != nil {
			_ = tx.Rollback()

			switch recoveredFrom.(type) {
			case error:
				err = recoveredFrom.(error)
			case string:
				err = errors.New(recoveredFrom.(string))
			default:
				err = fmt.Errorf("unknown panic: %v", recoveredFrom)
			}
		}
	}()

	if err := exec(r.newRepositoryTx(tx)); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *repository) exec() boil.ContextExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}
