package amqp

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"net"
	"net/url"
	"time"
)

type ProducerConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

type Producer struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

// TODO: Change params of exchange and queue; add comments

func NewProducer(config ProducerConfig) (*Producer, error) {
	var err error
	p := &Producer{}
	url := formUrl("amqp", config.Username, config.Password, config.Host, config.Port)

	p.conn, err = amqp091.Dial(url)
	if err != nil {
		return p, err
	}

	p.ch, err = p.conn.Channel()
	if err != nil {
		return p, err
	}

	err = p.ch.ExchangeDeclare("history-direct", "direct", true, false, false, false, nil)

	return p, err
}

func (p *Producer) Publish(body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.ch.PublishWithContext(ctx,
		"history-direct",
		"info",
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	return err
}

func (p *Producer) Shutdown() error {
	if err := p.ch.Close(); err != nil {
		return err
	}

	if err := p.conn.Close(); err != nil {
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
