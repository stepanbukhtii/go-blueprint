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

type company struct {
	exec boil.ContextExecutor
}

func NewCompany(exec boil.ContextExecutor) domain.CompanyRepository {
	return &company{exec: exec}
}

func (r *company) Add(ctx context.Context, company *domain.Company) error {
	m := convert.Company.Model(company)
	if err := m.Insert(ctx, r.exec, boil.Infer()); err != nil {
		return err
	}

	company.ID = m.ID
	company.CreatedAt = m.CreatedAt

	return nil
}

func (r *company) Update(ctx context.Context, company *domain.Company) error {
	m := convert.Company.Model(company)
	if _, err := m.Update(ctx, r.exec, boil.Infer()); err != nil {
		return err
	}
	company.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *company) Save(ctx context.Context, company *domain.Company) error {
	exists, err := r.Exists(ctx, model.CompanyWhere.ID.EQ(company.ID))
	if err != nil {
		return err
	}

	if exists {
		return r.Update(ctx, company)
	}

	return r.Add(ctx, company)
}

func (r *company) Find(ctx context.Context, mods ...qm.QueryMod) (domain.Company, error) {
	m, err := model.Companies(mods...).One(ctx, r.exec)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Company{}, domain.ErrCompanyNotFound
		}
		return domain.Company{}, err
	}
	return convert.Company.Domain(m), err
}

func (r *company) FindAll(ctx context.Context, mods ...qm.QueryMod) ([]domain.Company, error) {
	m, err := model.Companies(mods...).All(ctx, r.exec)
	if err != nil {
		return nil, err
	}
	return convert.Company.DomainSlice(m), err
}

func (r *company) FindAllPaginate(ctx context.Context, pagination api.Pagination, mods ...qm.QueryMod) ([]domain.Company, int64, error) {
	total, err := model.Companies(mods...).Count(ctx, r.exec)
	if err != nil {
		return nil, 0, err
	}

	mods = append(mods,
		qm.Limit(pagination.Limit()),
		qm.Offset(pagination.Offset()),
	)

	m, err := model.Companies(mods...).All(ctx, r.exec)
	if err != nil {
		return nil, 0, err
	}

	return convert.Company.DomainSlice(m), total, err
}

func (r *company) Exists(ctx context.Context, mods ...qm.QueryMod) (bool, error) {
	return model.Companies(mods...).Exists(ctx, r.exec)
}

func (r *company) Remove(ctx context.Context, id string) error {
	if _, err := model.Companies(model.CompanyWhere.ID.EQ(id)).DeleteAll(ctx, r.exec); err != nil {
		return err
	}
	return nil
}
