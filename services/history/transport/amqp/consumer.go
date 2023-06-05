package amqp

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"net"
	"net/url"
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
}

// TODO: Change params of exchange and queue; add comments

func NewConsumer(config ConsumerConfig) (*Consumer, error) {
	var err error
	c := &Consumer{}
	url := formUrl("amqp", config.Username, config.Password, config.Host, config.Port)

	c.conn, err = amqp091.Dial(url)
	if err != nil {
		return c, err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return c, err
	}

	err = c.ch.ExchangeDeclare("history-direct", "direct", true, false, false, false, nil)
	if err != nil {
		return c, err
	}

	q, err := c.ch.QueueDeclare(
		"",
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
		q.Name,
		"info",
		"history-direct",
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
		"",
		"",
		false,
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
			fmt.Printf("Received a message: %s\n", d.Body)
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
