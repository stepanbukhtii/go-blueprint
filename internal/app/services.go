package app

import (
	"fmt"

	"github.com/samber/do/v2"
	"github.com/stepanbukhtii/easy-tools/crypto"
	"github.com/stepanbukhtii/easy-tools/ejwt"
	"github.com/stepanbukhtii/easy-tools/grpc"
	"github.com/stepanbukhtii/go-blueprint/internal/clients/user"

	"github.com/stepanbukhtii/go-blueprint/internal/clients/randomuser"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/service"
	"github.com/stepanbukhtii/go-blueprint/internal/service/aggregator"
)

type Services struct {
	Auth     domain.AuthService
	Company  domain.CompanyService
	User     domain.UserService
	UserType domain.UserTypeService

	// Aggregators
	UserAggregator domain.UserAggregator
}

func (a *App) initServices() error {
	injector := do.New()

	do.ProvideValue(injector, a.Config)
	do.ProvideValue(injector, a.Repository)
	do.ProvideValue(injector, a.RabbitMQPublisher)
	do.ProvideValue(injector, a.KafkaProducer)
	do.ProvideValue(injector, a.NatsPublisher)

	if err := a.provideJWTGenerator(injector); err != nil {
		return err
	}

	if err := a.provideClients(injector); err != nil {
		return err
	}

	// aggregators
	do.Provide(injector, aggregator.NewUserAggregator)

	// services
	do.Provide(injector, service.NewAuth)
	do.Provide(injector, service.NewCompanyService)
	do.Provide(injector, service.NewUserService)
	do.Provide(injector, service.NewUserTypeService)

	a.Services = &Services{
		Auth:           do.MustInvoke[domain.AuthService](injector),
		Company:        do.MustInvoke[domain.CompanyService](injector),
		User:           do.MustInvoke[domain.UserService](injector),
		UserAggregator: do.MustInvoke[domain.UserAggregator](injector),
	}

	return nil
}

func (a *App) provideJWTGenerator(injector do.Injector) error {
	privateKey, err := crypto.ParseED25519PrivateKey(a.Config.JWT.PrivateKey)
	if err != nil {
		return err
	}

	jwtGenerator := ejwt.NewJWTGenerator(privateKey, a.Config.JWT.Issuer, a.Config.JWT.Audience, a.Config.JWT.ClaimsTTL)
	do.ProvideValue(injector, &jwtGenerator)

	return nil
}

func (a *App) provideClients(injector do.Injector) error {
	do.ProvideValue(injector, randomuser.NewClient(a.Config))

	grpcClient, err := grpc.NewClientConnection(fmt.Sprintf("localhost:%s", a.Config.GRPC.Port))
	if err != nil {
		return err
	}
	do.ProvideValue(injector, user.NewClient(grpcClient))

	return nil
}
