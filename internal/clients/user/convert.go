package user

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/pkg/grpc/proto"
)

func Domain(user *proto.User) domain.User {
	var lastRate float64
	if user.LastRate != nil {
		lastRate = *user.LastRate
	}

	var lockBalance *decimal.Decimal
	if user.LockBalance != nil {
		d := decimal.NewFromFloat(*user.LockBalance)
		lockBalance = &d
	}

	var lastLogin *time.Time
	if user.LastLogin != nil {
		t := user.LastLogin.AsTime()
		lastLogin = &t
	}

	return domain.User{
		ID:          user.Id,
		Name:        user.Name,
		Username:    user.Username,
		Password:    user.PublicName,
		PublicName:  user.PublicName,
		Description: user.Description,
		Age:         int(user.Age),
		InitialAge:  int(user.InitialAge),
		Rate:        user.Rate,
		LastRate:    lastRate,
		IsActive:    user.IsActive,
		ReadMessage: user.ReadMessage,
		Balance:     decimal.NewFromFloat(user.Balance),
		LockBalance: lockBalance,
		LastLogin:   lastLogin,
		CreatedAt:   user.CreatedAt.AsTime(),
		UpdatedAt:   user.UpdatedAt.AsTime(),
	}
}
