package nats

import (
	"github.com/stepanbukhtii/easy-tools/nats"
	"github.com/stepanbukhtii/go-blueprint/internal/app"
	"github.com/stepanbukhtii/go-blueprint/internal/service/events"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/nats/handlers"
)

type Subscriber struct {
	subscriber *nats.Subscriber
	services   *app.Services
}

func NewSubscriber(app *app.App) (*Subscriber, error) {
	natsSubscriber, err := nats.NewSubscriber(app.Config.NATS.ConnectionURI(), "queue_name")
	if err != nil {
		return nil, err
	}

	subscriber := &Subscriber{
		subscriber: natsSubscriber,
		services:   app.Services,
	}

	if err := subscriber.registerConsumers(); err != nil {
		return nil, err
	}

	return subscriber, nil
}

func (s *Subscriber) Shutdown() error {
	return s.subscriber.Shutdown()
}

func (s *Subscriber) registerConsumers() error {
	h := handlers.NewCompanyUpdate(s.services.Company, s.services.User)

	if err := s.subscriber.Subscribe(events.CompanyUpdatedEvent, h.Handle); err != nil {
		return err
	}

	return nil
}
