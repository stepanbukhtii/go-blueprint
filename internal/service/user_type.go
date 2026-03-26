package service

import (
	"context"
	"strings"

	"github.com/samber/do/v2"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository"
)

type userTypeService struct {
	repo repository.Repository
}

func NewUserTypeService(injector do.Injector) (domain.UserTypeService, error) {
	return &userTypeService{
		repo: do.MustInvoke[repository.Repository](injector),
	}, nil
}

func (s *userTypeService) List(ctx context.Context) ([]domain.UserType, error) {
	return s.repo.UserType().FindAll(ctx)
}

func (s *userTypeService) Get(ctx context.Context, id string) (domain.UserType, error) {
	return s.repo.UserType().Find(ctx, id)
}

func (s *userTypeService) Create(ctx context.Context, request domain.CreateUserTypeInput) (domain.UserType, error) {
	userType := domain.UserType{
		Code:    request.Code,
		Name:    request.Name,
		IsAdmin: request.IsAdmin,
	}

	if err := s.repo.UserType().Add(ctx, &userType); err != nil {
		return domain.UserType{}, err
	}

	return userType, nil
}

func (s *userTypeService) Update(ctx context.Context, request domain.UpdateUserTypeInput) (domain.UserType, error) {
	userType, err := s.repo.UserType().Find(ctx, request.Code)
	if err != nil {
		return domain.UserType{}, err
	}

	if request.Name != "" {
		userType.Name = strings.TrimSpace(request.Name)
	}

	if request.IsAdmin != nil {
		userType.IsAdmin = *request.IsAdmin
	}

	if err := s.repo.UserType().Update(ctx, &userType); err != nil {
		return domain.UserType{}, err
	}

	return userType, nil
}

func (s *userTypeService) Delete(ctx context.Context, id string) error {
	return s.repo.UserType().Remove(ctx, id)
}
