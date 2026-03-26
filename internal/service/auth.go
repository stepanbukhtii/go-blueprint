package service

import (
	"context"

	"github.com/samber/do/v2"
	"github.com/stepanbukhtii/easy-tools/ejwt"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/model"
	"golang.org/x/crypto/bcrypt"
)

type auth struct {
	repo         repository.Repository
	jwtGenerator *ejwt.JWTGenerator
}

func NewAuth(injector do.Injector) (domain.AuthService, error) {
	return &auth{
		repo:         do.MustInvoke[repository.Repository](injector),
		jwtGenerator: do.MustInvoke[*ejwt.JWTGenerator](injector),
	}, nil
}

func (s *auth) Login(ctx context.Context, request domain.LoginInput) (string, error) {
	user, err := s.repo.User().Find(ctx, model.UserWhere.Username.EQ(request.Username))
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return "", err
	}

	var token string
	if user.UserType.IsAdmin {
		token, err = s.jwtGenerator.GenerateTokenWthRoles(user.ID, []string{domain.RoleAdmin})
	} else {
		token, err = s.jwtGenerator.GenerateToken(user.ID)
	}
	if err != nil {
		return "", err
	}

	return token, nil
}
