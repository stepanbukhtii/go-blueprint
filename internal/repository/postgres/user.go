package postgres

import (
	"context"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/jackc/pgx/v5"
	"github.com/stepanbukhtii/easy-tools/rest/api"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/convert"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/model"
)

type user struct {
	exec boil.ContextExecutor
}

func NewUser(exec boil.ContextExecutor) domain.UserRepository {
	return &user{exec: exec}
}

func (r *user) Add(ctx context.Context, user *domain.User) error {
	m := convert.User.Model(user)
	if err := m.Insert(ctx, r.exec, boil.Infer()); err != nil {
		return err
	}

	user.ID = m.ID
	user.CreatedAt = m.CreatedAt

	return nil
}

func (r *user) Update(ctx context.Context, user *domain.User) error {
	m := convert.User.Model(user)
	if _, err := m.Update(ctx, r.exec, boil.Infer()); err != nil {
		return err
	}

	user.UpdatedAt = m.UpdatedAt

	return nil
}

func (r *user) Save(ctx context.Context, user *domain.User) error {
	exists, err := r.Exists(ctx, model.UserWhere.ID.EQ(user.ID))
	if err != nil {
		return err
	}

	if exists {
		return r.Update(ctx, user)
	}

	return r.Add(ctx, user)
}

func (r *user) Find(ctx context.Context, mods ...qm.QueryMod) (domain.User, error) {
	m, err := model.Users(mods...).One(ctx, r.exec)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return convert.User.Domain(m), err
}

func (r *user) FindAll(ctx context.Context, mods ...qm.QueryMod) ([]domain.User, error) {
	m, err := model.Users(mods...).All(ctx, r.exec)
	if err != nil {
		return nil, err
	}
	return convert.User.DomainSlice(m), err
}

func (r *user) FindAllPaginate(ctx context.Context, pagination api.Pagination, mods ...qm.QueryMod) ([]domain.User, int64, error) {
	total, err := model.Users(mods...).Count(ctx, r.exec)
	if err != nil {
		return nil, 0, err
	}

	mods = append(mods,
		qm.Limit(pagination.Limit()),
		qm.Offset(pagination.Offset()),
	)

	m, err := model.Users(mods...).All(ctx, r.exec)
	if err != nil {
		return nil, 0, err
	}

	return convert.User.DomainSlice(m), total, err
}

func (r *user) Exists(ctx context.Context, mods ...qm.QueryMod) (bool, error) {
	return model.Users(mods...).Exists(ctx, r.exec)
}

func (r *user) Remove(ctx context.Context, id string) error {
	if _, err := model.Users(model.UserWhere.ID.EQ(id)).DeleteAll(ctx, r.exec); err != nil {
		return err
	}
	return nil
}
