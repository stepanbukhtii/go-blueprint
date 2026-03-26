package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/samber/do/v2"
	"github.com/stepanbukhtii/easy-tools/kafka"
	"github.com/stepanbukhtii/easy-tools/rabbitmq"
	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stepanbukhtii/go-blueprint/internal/service/events"
	"golang.org/x/crypto/bcrypt"

	"github.com/stepanbukhtii/go-blueprint/internal/clients/randomuser"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/repository"
	"github.com/stepanbukhtii/go-blueprint/internal/repository/postgres/model"
)

type userService struct {
	repo              repository.Repository
	rabbitmqPublisher *rabbitmq.Publisher
	kafkaProducer     *kafka.Producer
	userAggregator    domain.UserAggregator
	randomUserClient  randomuser.Client
}

func NewUserService(injector do.Injector) (domain.UserService, error) {
	return &userService{
		repo:              do.MustInvoke[repository.Repository](injector),
		rabbitmqPublisher: do.MustInvoke[*rabbitmq.Publisher](injector),
		kafkaProducer:     do.MustInvoke[*kafka.Producer](injector),
		userAggregator:    do.MustInvoke[domain.UserAggregator](injector),
		randomUserClient:  do.MustInvoke[randomuser.Client](injector),
	}, nil
}

func (s *userService) ListPaginate(ctx context.Context, request api.Pagination) ([]domain.User, int64, error) {
	return s.repo.User().FindAllPaginate(ctx, request)
}

func (s *userService) Get(ctx context.Context, id string) (domain.User, error) {
	return s.userAggregator.Get(ctx, id)
}

func (s *userService) Create(ctx context.Context, request domain.CreateUserInput) (domain.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		Name:     request.Name,
		Username: request.Username,
		Password: string(passwordHash),
		UserType: domain.UserType{Code: domain.UserTypeDefault},
	}

	if err := s.repo.User().Add(ctx, &user); err != nil {
		return domain.User{}, err
	}

	slog.With(
		slog.String("user.id", user.ID),
		slog.String("user.username", user.Username),
	).InfoContext(ctx, "user created")

	eventData := events.NewEventUserCreatedData(user)
	if err := s.kafkaProducer.Produce(ctx, events.UserCreatedEvent, eventData); err != nil {
		return domain.User{}, err
	}

	slog.With(
		slog.String("user.id", user.ID),
		slog.String("event.name", events.UserCreatedEvent),
	).InfoContext(ctx, "event published")

	return user, nil
}

func (s *userService) Update(ctx context.Context, request domain.UpdateUserInput) (domain.User, error) {
	user, err := s.repo.User().Find(ctx, model.UserWhere.ID.EQ(request.UserID))
	if err != nil {
		return domain.User{}, err
	}

	if request.Name != "" {
		user.Name = strings.TrimSpace(request.Name)
	}

	if request.Username != "" {
		user.Username = strings.ToLower(strings.TrimSpace(request.Username))
	}

	if err := s.repo.User().Update(ctx, &user); err != nil {
		return domain.User{}, err
	}

	slog.With(
		slog.String("user.id", user.ID),
		slog.String("user.username", user.Username),
	).InfoContext(ctx, "user updated")

	eventData := events.NewEventUserUpdatedData(user)
	if err := s.rabbitmqPublisher.PublishQueue(ctx, events.UserUpdatedEvent, eventData); err != nil {
		return domain.User{}, err
	}

	slog.With(
		slog.String("user.id", user.ID),
		slog.String("event.name", events.UserUpdatedEvent),
	).InfoContext(ctx, "event was published")

	return user, nil
}

func (s *userService) GeneratePublicName(ctx context.Context, id string) error {
	user, err := s.repo.User().Find(ctx, model.UserWhere.ID.EQ(id))
	if err != nil {
		return err
	}

	randomUser, err := s.randomUserClient.GetRandomUser(ctx)
	if err != nil {
		return err
	}

	user.PublicName = fmt.Sprintf("%s %s", randomUser.Name.First, randomUser.Name.Last)

	if err := s.repo.User().Update(ctx, &user); err != nil {
		return err
	}

	slog.With(
		slog.String("user.id", id),
		slog.String("user.public_name", user.PublicName),
	).InfoContext(ctx, "public name generated")

	return nil
}

func (s *userService) UpdateRate(ctx context.Context, id string, lastRate float64) error {
	user, err := s.repo.User().Find(ctx, model.UserWhere.ID.EQ(id))
	if err != nil {
		return err
	}

	user.LastRate = lastRate
	if user.Rate >= 0 {
		user.Rate = 0
	}

	if err := s.repo.User().Update(ctx, &user); err != nil {
		return err
	}

	slog.With(
		slog.String("user.id", id),
		slog.String("user.public_name", user.PublicName),
	).InfoContext(ctx, "rate updated")

	return nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.User().Remove(ctx, id)
}
