package handlers

import (
	"context"
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stepanbukhtii/easy-tools/rabbitmq"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
)

type UserUpdate struct {
	service domain.UserService
}

func NewUserUpdate(service domain.UserService) *UserUpdate {
	return &UserUpdate{
		service: service,
	}
}

func DeclareUserUpdate(declare rabbitmq.Declare) error {
	if err := declare.Queue(domain.UserUpdatedEvent); err != nil {
		return err
	}

	return nil
}

func (u UserUpdate) Handle(ctx context.Context, msg amqp.Delivery) error {
	var request domain.EventUserUpdatedData
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}

	slog.With(
		slog.String("user.id", request.UserID),
		slog.String("event.name", domain.UserUpdatedEvent),
	).InfoContext(ctx, "rabbitMQ event received")

	if err := u.service.UpdateRate(ctx, request.UserID, request.Rate); err != nil {
		return err
	}

	return nil
}
