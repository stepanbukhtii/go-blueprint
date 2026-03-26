package postgres

import (
	"context"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/jackc/pgx/v5"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/convert"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/model"
)

type userType struct {
	exec boil.ContextExecutor
}

func NewUserType(exec boil.ContextExecutor) domain.UserTypeRepository {
	return &userType{exec: exec}
}

func (r *userType) Add(ctx context.Context, userType *domain.UserType) error {
	m := convert.UserType.Model(userType)
	if err := m.Insert(ctx, r.exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (r *userType) Update(ctx context.Context, userType *domain.UserType) error {
	m := convert.UserType.Model(userType)
	if _, err := m.Update(ctx, r.exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (r *userType) Find(ctx context.Context, code string) (domain.UserType, error) {
	m, err := model.UserTypes(model.UserTypeWhere.Code.EQ(code)).One(ctx, r.exec)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserType{}, domain.ErrUserTypeNotFound
		}
		return domain.UserType{}, err
	}
	return convert.UserType.Domain(m), err
}

func (r *userType) FindAll(ctx context.Context) ([]domain.UserType, error) {
	m, err := model.UserTypes().All(ctx, r.exec)
	if err != nil {
		return nil, err
	}
	return convert.UserType.DomainSlice(m), err
}

func (r *userType) Remove(ctx context.Context, code string) error {
	if _, err := model.UserTypes(model.UserTypeWhere.Code.EQ(code)).DeleteAll(ctx, r.exec); err != nil {
		return err
	}
	return nil
}
