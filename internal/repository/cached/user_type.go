package cached

import (
	"context"

	"github.com/stepanbukhtii/easy-tools/cache"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
)

type userTypeRepository struct {
	cache cache.MapCache[domain.UserType]
	repo  domain.UserTypeRepository
}

func NewUserTypeRepository(cache cache.MapCache[domain.UserType], repo domain.UserTypeRepository) domain.UserTypeRepository {
	return &userTypeRepository{
		cache: cache,
		repo:  repo,
	}
}

func (r *userTypeRepository) Add(ctx context.Context, userType *domain.UserType) error {
	if err := r.cache.Set(ctx, userType.Code, *userType); err != nil {
		return err
	}
	return r.repo.Add(ctx, userType)
}

func (r *userTypeRepository) Update(ctx context.Context, userType *domain.UserType) error {
	if err := r.cache.Set(ctx, userType.Code, *userType); err != nil {
		return err
	}
	return r.repo.Update(ctx, userType)
}

func (r *userTypeRepository) Find(ctx context.Context, code string) (domain.UserType, error) {
	userType, err := r.cache.Get(ctx, code)
	if err == nil {
		return userType, nil
	}

	userType, err = r.repo.Find(ctx, code)
	if err != nil {
		return domain.UserType{}, err
	}

	_ = r.refreshCache(ctx)

	return userType, nil
}

func (r *userTypeRepository) FindAll(ctx context.Context) ([]domain.UserType, error) {
	userTypes, err := r.cache.GetAll(ctx)
	if err == nil && len(userTypes) > 0 {
		return userTypes, nil
	}

	userTypes, err = r.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	_ = r.cache.SetAll(ctx, userTypes, func(v domain.UserType) string { return v.Code })

	return userTypes, nil
}

func (r *userTypeRepository) Remove(ctx context.Context, code string) error {
	if err := r.cache.Delete(ctx, code); err != nil {
		return err
	}
	return r.repo.Remove(ctx, code)
}

func (r *userTypeRepository) refreshCache(ctx context.Context) error {
	userTypes, err := r.repo.FindAll(ctx)
	if err != nil {
		return err
	}

	if err = r.cache.DeleteAll(ctx); err != nil {
		return err
	}

	return r.cache.SetAll(ctx, userTypes, func(v domain.UserType) string { return v.Code })
}
