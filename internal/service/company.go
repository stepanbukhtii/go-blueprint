package service

import (
	"context"
	"strings"

	"github.com/samber/do/v2"
	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/rest/api"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/model"
)

type companyService struct {
	repo repository.Repository
}

func NewCompanyService(injector do.Injector) (domain.CompanyService, error) {
	return &companyService{
		repo: do.MustInvoke[repository.Repository](injector),
	}, nil
}

func (s *companyService) ListPaginate(ctx context.Context, request api.Pagination) ([]domain.Company, int64, error) {
	return s.repo.Company().FindAllPaginate(ctx, request)
}

func (s *companyService) Get(ctx context.Context, id string) (domain.Company, error) {
	return s.repo.Company().Find(ctx, model.CompanyWhere.ID.EQ(id))
}

func (s *companyService) GetCompanyByOwner(ctx context.Context) ([]domain.Company, error) {
	clientInfo := econtext.ClientInfo(ctx)
	if clientInfo.Subject == "" {
		return []domain.Company{}, domain.ErrCompanyNotFound
	}
	return s.repo.Company().FindAll(ctx, model.CompanyWhere.OwnerID.EQ(clientInfo.Subject))
}

func (s *companyService) Create(ctx context.Context, request domain.CreateCompanyInput) (domain.Company, error) {
	clientInfo := econtext.ClientInfo(ctx)

	company := domain.Company{
		Name:      request.Name,
		OwnerID:   clientInfo.Subject,
		ManagerID: request.ManagerID,
		LogoURL:   request.LogoURL,
	}

	if err := s.repo.Company().Add(ctx, &company); err != nil {
		return domain.Company{}, err
	}

	return company, nil
}

func (s *companyService) CreateMultiple(ctx context.Context, request []domain.CreateCompanyInput) ([]domain.Company, error) {
	var companies []domain.Company

	for _, req := range request {
		company := domain.Company{
			Name:      req.Name,
			OwnerID:   req.OwnerID,
			ManagerID: req.ManagerID,
			LogoURL:   req.LogoURL,
		}
		if err := s.repo.Company().Add(ctx, &company); err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

func (s *companyService) Update(ctx context.Context, request domain.UpdateCompanyInput) (domain.Company, error) {
	company, err := s.repo.Company().Find(ctx, model.CompanyWhere.ID.EQ(request.CompanyID))
	if err != nil {
		return domain.Company{}, err
	}

	if request.Name != "" {
		company.Name = strings.TrimSpace(request.Name)
	}

	if request.LogoURL != "" {
		company.LogoURL = request.LogoURL
	}

	if err := s.repo.Company().Update(ctx, &company); err != nil {
		return domain.Company{}, err
	}

	return company, nil
}

func (s *companyService) Delete(ctx context.Context, id string) error {
	return s.repo.Company().Remove(ctx, id)
}
