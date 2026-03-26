package handlers

import (
	"context"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/grpc/handlers/response"
	"github.com/stepanbukhtii/go-blueprint/pkg/grpc/proto"
)

type UserHandler struct {
	proto.UnimplementedUserServiceServer

	userAggregator domain.UserAggregator
}

func NewUserHandler(userAggregator domain.UserAggregator) *UserHandler {
	return &UserHandler{
		userAggregator: userAggregator,
	}
}

func (s *UserHandler) One(ctx context.Context, request *proto.OneRequest) (*proto.User, error) {
	user, err := s.userAggregator.Get(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return response.NewUser(user), nil
}
