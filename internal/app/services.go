package app

import (
	"github.com/samber/do/v2"
	"github.com/stepanbukhtii/easy-tools/crypto"
	"github.com/stepanbukhtii/easy-tools/ejwt"
	"github.com/stepanbukhtii/easy-tools/kafka"
	"github.com/stepanbukhtii/easy-tools/rabbitmq"

	"github.com/stepanbukhtii/go-blueprint/internal/config"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/external/randomuser"
	"github.com/stepanbukhtii/go-blueprint/internal/repository"
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

func NewServices(
	cfg config.Config,
	repo repository.Repository,
	rabbitMQPublisher *rabbitmq.Publisher,
	kafkaProducer *kafka.Producer,
) (*Services, error) {
	injector := do.New()

	do.ProvideValue(injector, cfg)
	do.ProvideValue(injector, repo)
	do.ProvideValue(injector, rabbitMQPublisher)
	do.ProvideValue(injector, kafkaProducer)

	// jwt generator
	privateKey, err := crypto.ParseED25519PrivateKey(cfg.JWT.PrivateKey)
	if err != nil {
		return nil, err
	}
	jwtGenerator := ejwt.NewJWTGenerator(privateKey, cfg.JWT.Issuer, cfg.JWT.Audience, cfg.JWT.ClaimsTTL)
	do.ProvideValue(injector, &jwtGenerator)

	// externals
	do.Provide(injector, randomuser.NewClient)

	// aggregators
	do.Provide(injector, aggregator.NewUserAggregator)

	// services
	do.Provide(injector, service.NewAuth)
	do.Provide(injector, service.NewCompanyService)
	do.Provide(injector, service.NewUserService)
	do.Provide(injector, service.NewUserTypeService)

	return &Services{
		Auth:           do.MustInvoke[domain.AuthService](injector),
		Company:        do.MustInvoke[domain.CompanyService](injector),
		User:           do.MustInvoke[domain.UserService](injector),
		UserAggregator: do.MustInvoke[domain.UserAggregator](injector),
	}, nil
}
