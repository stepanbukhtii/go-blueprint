package handlers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/twmb/franz-go/pkg/kgo"
)

type UserCreated struct {
	service domain.UserService
}

func NewUserCreated(service domain.UserService) *UserCreated {
	return &UserCreated{
		service: service,
	}
}

func (u *UserCreated) Handle(ctx context.Context, record *kgo.Record) error {
	var request domain.EventUserCreatedData
	if err := json.Unmarshal(record.Value, &request); err != nil {
		return err
	}

	slog.With(
		slog.String("user.id", request.UserID),
		slog.String("event.name", domain.UserCreatedEvent),
	).InfoContext(ctx, "kafka event received")

	if err := u.service.GeneratePublicName(ctx, request.UserID); err != nil {
		return err
	}

	return nil
}
