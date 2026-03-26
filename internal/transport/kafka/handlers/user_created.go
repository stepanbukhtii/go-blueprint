package handlers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/service/events"
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

func (h *UserCreated) Handle(ctx context.Context, record *kgo.Record) error {
	var request events.EventUserCreatedData
	if err := json.Unmarshal(record.Value, &request); err != nil {
		return err
	}

	slog.With(
		slog.String("user.id", request.UserID),
		slog.String("event.name", events.UserCreatedEvent),
	).InfoContext(ctx, "start generate public name")

	if err := h.service.GeneratePublicName(ctx, request.UserID); err != nil {
		return err
	}

	return nil
}
