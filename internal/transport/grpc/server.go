package grpc

import (
	"github.com/stepanbukhtii/easy-tools/grpc"

	"github.com/stepanbukhtii/go-blueprint/internal/app"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/grpc/handlers"
	pkgproto "github.com/stepanbukhtii/go-blueprint/pkg/grpc/proto"
)

func NewServer(app *app.App) grpc.Server {
	s := grpc.NewServer()

	pkgproto.RegisterUserServiceServer(s, handlers.NewUserHandler(app.Services.User))

	return s
}
