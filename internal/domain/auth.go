package domain

import (
	"context"
)

const RoleAdmin = "admin"

type LoginInput struct {
	Username string
	Password string
}

type AuthService interface {
	Login(ctx context.Context, request LoginInput) (string, error)
}
