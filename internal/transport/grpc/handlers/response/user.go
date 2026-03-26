package response

import (
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/pkg/grpc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewUser(user domain.User) *proto.User {
	var lockBalance *float64
	if user.LockBalance != nil {
		f := user.LockBalance.InexactFloat64()
		lockBalance = &f
	}

	var lastLogin *timestamppb.Timestamp
	if user.LastLogin != nil {
		lastLogin = timestamppb.New(*user.LastLogin)
	}

	return &proto.User{
		Id:          user.ID,
		Name:        user.Name,
		Username:    user.Username,
		PublicName:  user.PublicName,
		Description: user.Description,
		Age:         int64(user.Age),
		InitialAge:  int64(user.InitialAge),
		Rate:        user.Rate,
		LastRate:    &user.LastRate,
		IsActive:    user.IsActive,
		ReadMessage: user.ReadMessage,
		Balance:     user.Balance.InexactFloat64(),
		LockBalance: lockBalance,
		LastLogin:   lastLogin,
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}
}
