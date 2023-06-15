package amqp

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rvinnie/lightstream/services/history/service"
	"net"
	"net/url"
)

const (
	exchangeName = "history"
	exchangeKind = "fanout"
	queueName    = ""
)

type ConsumerConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

type Consumer struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel

	notificationsService service.Notifications
}

func NewConsumer(config ConsumerConfig, notificationsService service.Notifications) (*Consumer, error) {
	var err error
	c := &Consumer{}

	c.notificationsService = notificationsService

	url := formUrl("amqp", config.Username, config.Password, config.Host, config.Port)

	c.conn, err = amqp091.Dial(url)
	if err != nil {
		return c, err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return c, err
	}

	err = c.ch.ExchangeDeclare(exchangeName, exchangeKind, true, false, false, false, nil)
	if err != nil {
		return c, err
	}

	_, err = c.ch.QueueDeclare(
		queueName,
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return c, err
	}

	err = c.ch.QueueBind(
		queueName,
		"",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return c, err
	}

	return c, err
}

func (c *Consumer) Consume() error {
	deliveries, err := c.ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range deliveries {
			c.notificationsService.Create(context.Background(), string(d.Body))
		}
	}()

	return nil
}

func (c *Consumer) Shutdown() error {
	if err := c.ch.Close(); err != nil {
		return err
	}

	if err := c.conn.Close(); err != nil {
		return err
	}

	return nil
}

func formUrl(scheme, username, password, host, port string) string {
	var u = url.URL{
		Scheme: scheme,
		User:   url.UserPassword(username, password),
		Host:   net.JoinHostPort(host, port),
	}
	return u.String()
}
