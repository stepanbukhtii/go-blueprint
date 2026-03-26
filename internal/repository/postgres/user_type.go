package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/stephenafamo/bob"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/convert"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/models"
)

type userType struct {
	exec bob.Executor
}

func NewUserType(exec bob.Executor) domain.UserTypeRepository {
	return &userType{exec: exec}
}

func (r *userType) Add(ctx context.Context, userType *domain.UserType) error {
	_, err := models.UserTypes.Insert(convert.UserType.Setter(userType)).One(ctx, r.exec)
	return err
}

func (r *userType) Update(ctx context.Context, userType *domain.UserType) error {
	_, err := models.UserTypes.Update(
		models.UpdateWhere.UserTypes.Code.EQ(userType.Code),
		convert.UserType.Setter(userType).UpdateMod(),
	).Exec(ctx, r.exec)
	return err
}

func (r *userType) Find(ctx context.Context, code string) (domain.UserType, error) {
	m, err := models.UserTypes.Query(models.SelectWhere.UserTypes.Code.EQ(code)).One(ctx, r.exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.UserType{}, domain.ErrUserTypeNotFound
		}
		return domain.UserType{}, err
	}
	return convert.UserType.Domain(m), nil
}

func (r *userType) FindAll(ctx context.Context) ([]domain.UserType, error) {
	m, err := models.UserTypes.Query().All(ctx, r.exec)
	if err != nil {
		return nil, err
	}
	return convert.UserType.DomainSlice(m), nil
}

func (r *userType) Remove(ctx context.Context, code string) error {
	_, err := models.UserTypes.Delete(models.DeleteWhere.UserTypes.Code.EQ(code)).Exec(ctx, r.exec)
	return err
}
