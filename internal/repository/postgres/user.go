package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/convert"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/dberrors"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/models"
)

type user struct {
	exec bob.Executor
}

func NewUser(exec bob.Executor) domain.UserRepository {
	return &user{exec: exec}
}

func (r *user) Add(ctx context.Context, u *domain.User) error {
	m, err := models.Users.Insert(convert.User.Setter(u)).One(ctx, r.exec)
	if err != nil {
		if errors.Is(dberrors.UserErrors.ErrUniqueUserUsernameUnique, err) {
			return domain.ErrUsernameIsExists
		}
		return err
	}
	u.ID = m.ID
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *user) Update(ctx context.Context, u *domain.User) error {
	m, err := models.Users.Update(
		models.UpdateWhere.Users.ID.EQ(u.ID),
		convert.User.Setter(u).UpdateMod(),
	).One(ctx, r.exec)
	if err != nil {
		return err
	}
	u.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *user) Save(ctx context.Context, u *domain.User) error {
	exists, err := r.Exists(ctx, models.SelectWhere.Users.ID.EQ(u.ID))
	if err != nil {
		return err
	}
	if exists {
		return r.Update(ctx, u)
	}
	return r.Add(ctx, u)
}

func (r *user) Find(ctx context.Context, mods ...bob.Mod[*dialect.SelectQuery]) (domain.User, error) {
	m, err := models.Users.Query(mods...).One(ctx, r.exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return convert.User.Domain(m), nil
}

func (r *user) FindAll(ctx context.Context, mods ...bob.Mod[*dialect.SelectQuery]) ([]domain.User, error) {
	m, err := models.Users.Query(mods...).All(ctx, r.exec)
	if err != nil {
		return nil, err
	}
	return convert.User.DomainSlice(m), nil
}

func (r *user) FindAllPaginate(
	ctx context.Context,
	pagination api.Pagination,
	mods ...bob.Mod[*dialect.SelectQuery],
) ([]domain.User, int64, error) {
	total, err := models.Users.Query(mods...).Count(ctx, r.exec)
	if err != nil {
		return nil, 0, err
	}

	mods = append(
		mods,
		sm.Limit(pagination.Limit()),
		sm.Offset(pagination.Offset()),
	)

	m, err := models.Users.Query(mods...).All(ctx, r.exec)
	if err != nil {
		return nil, 0, err
	}

	return convert.User.DomainSlice(m), total, nil
}

func (r *user) Exists(ctx context.Context, mods ...bob.Mod[*dialect.SelectQuery]) (bool, error) {
	return models.Users.Query(mods...).Exists(ctx, r.exec)
}

func (r *user) Remove(ctx context.Context, id string) error {
	_, err := models.Users.Delete(models.DeleteWhere.Users.ID.EQ(id)).Exec(ctx, r.exec)
	return err
}
