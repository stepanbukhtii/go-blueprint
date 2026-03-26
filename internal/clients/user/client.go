package user

import (
	"context"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	pkgproto "github.com/stepanbukhtii/go-blueprint/pkg/grpc/proto"
	"google.golang.org/grpc"
)

type Client interface {
	One(ctx context.Context, id string) (domain.User, error)
}

type client struct {
	userGRPCClient pkgproto.UserServiceClient
}

func NewClient(conn *grpc.ClientConn) Client {
	return &client{
		userGRPCClient: pkgproto.NewUserServiceClient(conn),
	}
}

func (s *client) One(ctx context.Context, id string) (domain.User, error) {
	user, err := s.userGRPCClient.One(ctx, &pkgproto.OneRequest{Id: id})
	if err != nil {
		return domain.User{}, err
	}
	return Domain(user), nil
}
