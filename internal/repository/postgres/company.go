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
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/models"
)

type company struct {
	exec bob.Executor
}

func NewCompany(exec bob.Executor) domain.CompanyRepository {
	return &company{exec: exec}
}

func (r *company) Add(ctx context.Context, company *domain.Company) error {
	m, err := models.Companies.Insert(convert.Company.Setter(company)).One(ctx, r.exec)
	if err != nil {
		return err
	}
	company.ID = m.ID
	company.CreatedAt = m.CreatedAt
	company.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *company) Update(ctx context.Context, company *domain.Company) error {
	m, err := models.Companies.Update(
		models.UpdateWhere.Companies.ID.EQ(company.ID),
		convert.Company.Setter(company).UpdateMod(),
	).One(ctx, r.exec)
	if err != nil {
		return err
	}
	company.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *company) Save(ctx context.Context, company *domain.Company) error {
	exists, err := r.Exists(ctx, models.SelectWhere.Companies.ID.EQ(company.ID))
	if err != nil {
		return err
	}
	if exists {
		return r.Update(ctx, company)
	}
	return r.Add(ctx, company)
}

func (r *company) Find(ctx context.Context, mods ...bob.Mod[*dialect.SelectQuery]) (domain.Company, error) {
	m, err := models.Companies.Query(mods...).One(ctx, r.exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Company{}, domain.ErrCompanyNotFound
		}
		return domain.Company{}, err
	}
	return convert.Company.Domain(m), nil
}

func (r *company) FindAll(ctx context.Context, mods ...bob.Mod[*dialect.SelectQuery]) ([]domain.Company, error) {
	m, err := models.Companies.Query(mods...).All(ctx, r.exec)
	if err != nil {
		return nil, err
	}
	return convert.Company.DomainSlice(m), nil
}

func (r *company) FindAllPaginate(
	ctx context.Context,
	pagination api.Pagination,
	mods ...bob.Mod[*dialect.SelectQuery],
) ([]domain.Company, int64, error) {
	total, err := models.Companies.Query(mods...).Count(ctx, r.exec)
	if err != nil {
		return nil, 0, err
	}

	mods = append(
		mods,
		sm.Limit(pagination.Limit()),
		sm.Offset(pagination.Offset()),
	)

	m, err := models.Companies.Query(mods...).All(ctx, r.exec)
	if err != nil {
		return nil, 0, err
	}

	return convert.Company.DomainSlice(m), total, nil
}

func (r *company) Exists(ctx context.Context, mods ...bob.Mod[*dialect.SelectQuery]) (bool, error) {
	return models.Companies.Query(mods...).Exists(ctx, r.exec)
}

func (r *company) Remove(ctx context.Context, id string) error {
	_, err := models.Companies.Delete(models.DeleteWhere.Users.ID.EQ(id)).Exec(ctx, r.exec)
	return err
}
