package request

import "github.com/stepanbukhtii/go-blueprint/internal/domain"

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r Login) ToDomain() domain.LoginInput {
	return domain.LoginInput{
		Username: r.Username,
		Password: r.Password,
	}
}
