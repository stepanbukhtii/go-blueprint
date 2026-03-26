package response

import (
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
)

type User struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Username    string           `json:"username"`
	Password    string           `json:"password"`
	PublicName  string           `json:"public_name,omitempty"`
	Description string           `json:"description"`
	Age         int              `json:"age"`
	InitialAge  int              `json:"initial_age,omitempty"`
	Rate        float64          `json:"rate"`
	LastRate    float64          `json:"last_rate,omitempty"`
	IsActive    bool             `json:"is_active"`
	ReadMessage *bool            `json:"read_message,omitempty"`
	Balance     decimal.Decimal  `json:"balance"`
	LockBalance *decimal.Decimal `json:"lock_balance"`
	LastLogin   *time.Time       `json:"last_login"`
}

func NewUser(u domain.User) User {
	return User{
		ID:          u.ID,
		Name:        u.Name,
		Username:    u.Username,
		Password:    u.Password,
		PublicName:  u.PublicName,
		Description: u.Description,
		Age:         u.Age,
		InitialAge:  u.InitialAge,
		Rate:        u.Rate,
		LastRate:    u.LastRate,
		IsActive:    u.IsActive,
		ReadMessage: u.ReadMessage,
		Balance:     u.Balance,
		LockBalance: u.LockBalance,
		LastLogin:   u.LastLogin,
	}
}

func NewUsers(users []domain.User) []User {
	return lo.Map(users, func(u domain.User, index int) User { return NewUser(u) })
}
