package request

import "github.com/stepanbukhtii/go-blueprint/internal/domain"

type UserURI struct {
	UserID string `uri:"user_id" binding:"required,uuid" swaggerignore:"true"`
}

type CreateUser struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r CreateUser) ToDomain() domain.CreateUserInput {
	return domain.CreateUserInput{
		Name:     r.Name,
		Username: r.Username,
		Password: r.Password,
	}
}

type UpdateUser struct {
	UserURI
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (r UpdateUser) ToDomain() domain.UpdateUserInput {
	return domain.UpdateUserInput{
		UserID:   r.UserID,
		Name:     r.Name,
		Username: r.Username,
	}
}
