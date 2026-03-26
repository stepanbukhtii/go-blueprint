package rabbitmq

import (
	"github.com/stepanbukhtii/easy-tools/rabbitmq"
	"github.com/stepanbukhtii/go-blueprint/internal/app"
	"github.com/stepanbukhtii/go-blueprint/internal/service/events"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/rabbitmq/handlers"
)

type Consumers struct {
	group              *rabbitmq.ConsumerGroup
	rabbitMQConnection *rabbitmq.Connection
	services           *app.Services
}

func NewConsumers(app *app.App) (*rabbitmq.ConsumerGroup, error) {
	rabbitMQConnection, err := rabbitmq.NewConnection(app.Config.RabbitMQ.ConnectionURI())
	if err != nil {
		return nil, err
	}

	consumers := &Consumers{
		group:              rabbitmq.NewConsumerWithConnection(rabbitMQConnection, app.Config.Service.Name),
		rabbitMQConnection: rabbitMQConnection,
		services:           app.Services,
	}

	if err := consumers.declareConsumers(); err != nil {
		return nil, err
	}

	consumers.registerConsumers()

	return consumers.group, nil
}

func (c *Consumers) declareConsumers() error {
	declare := c.rabbitMQConnection.Declare()

	if err := handlers.DeclareUserUpdate(declare); err != nil {
		return err
	}

	return nil
}

func (c *Consumers) registerConsumers() {
	h := handlers.NewUserUpdate(c.services.User)

	c.group.Add(events.UserUpdatedEvent, h.Handle)
}
