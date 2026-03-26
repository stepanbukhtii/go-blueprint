package aggregator

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/stepanbukhtii/easy-tools/rest/api"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/model"
)

type userAggregator struct {
	repo repository.Repository
}

func NewUserAggregator(injector do.Injector) (domain.UserAggregator, error) {
	return &userAggregator{
		repo: do.MustInvoke[repository.Repository](injector),
	}, nil
}

func (a *userAggregator) ListPaginate(ctx context.Context, request api.Pagination) ([]domain.User, int64, error) {
	users, total, err := a.repo.User().FindAllPaginate(ctx, request, qm.Load(model.UserRels.ManagerCompanies))
	if err != nil {
		return nil, 0, err
	}

	if err := a.setManyRelations(ctx, users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (a *userAggregator) Get(ctx context.Context, id string) (domain.User, error) {
	user, err := a.repo.User().Find(
		ctx,
		model.UserWhere.ID.EQ(id),
		qm.Load(model.UserRels.ManagerCompanies),
	)
	if err != nil {
		return domain.User{}, err
	}

	if err := a.setOneRelations(ctx, &user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (a *userAggregator) setOneRelations(ctx context.Context, user *domain.User) error {
	var err error

	user.UserType, err = a.repo.UserType().Find(ctx, user.UserType.Code)
	if err != nil {
		return err
	}

	return nil
}

func (a *userAggregator) setManyRelations(ctx context.Context, users []domain.User) error {
	if err := a.setUserTypes(ctx, users); err != nil {
		return err
	}

	return nil
}

func (a *userAggregator) setUserTypes(ctx context.Context, users []domain.User) error {
	userTypes, err := a.repo.UserType().FindAll(ctx)
	if err != nil {
		return err
	}

	userTypesMap := lo.KeyBy(userTypes, func(ut domain.UserType) string { return ut.Code })

	for i, _ := range users {
		users[i].UserType = userTypesMap[users[i].UserType.Code]
	}

	return nil
}
