package handlers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/service/events"
)

type CompanyUpdate struct {
	service     domain.CompanyService
	userService domain.UserService
}

func NewCompanyUpdate(service domain.CompanyService, userService domain.UserService) *CompanyUpdate {
	return &CompanyUpdate{
		service:     service,
		userService: userService,
	}
}

func (h CompanyUpdate) Handle(ctx context.Context, msg *nats.Msg) error {
	var request events.EventCompanyUpdatedData
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		return err
	}

	company, err := h.service.Get(ctx, request.CompanyID)
	if err != nil {
		return err
	}

	slog.With(
		slog.String("company.id", request.CompanyID),
		slog.String("event.name", events.UserUpdatedEvent),
	).InfoContext(ctx, "start update company owner rate")

	if company.OwnerID == request.OwnerID {
		if err := h.userService.UpdateRate(ctx, request.OwnerID, 100); err != nil {
			return err
		}
	}

	return nil
}
