package kafka

import (
	"github.com/stepanbukhtii/easy-tools/kafka"
	"github.com/stepanbukhtii/go-blueprint/internal/app"
	"github.com/stepanbukhtii/go-blueprint/internal/service/events"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/kafka/handlers"
)

type Consumers struct {
	group    *kafka.ConsumerGroup
	services *app.Services
}

func NewConsumers(app *app.App) *kafka.ConsumerGroup {
	consumers := &Consumers{
		group:    kafka.NewConsumerGroup(app.Config.Service.Name, app.Config.Kafka.Brokers...),
		services: app.Services,
	}

	consumers.registerConsumers()

	return consumers.group
}

func (c *Consumers) registerConsumers() {
	h := handlers.NewUserCreated(c.services.User)

	c.group.Add(events.UserCreatedEvent, h.Handle)
}
