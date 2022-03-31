package rabbitmq

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type rabbitMQChannel struct {
	uuid       string
	connection *amqp.Connection
	channel    *amqp.Channel
	confirms   chan amqp.Confirmation
}

func newRabbitChannel(conn *amqp.Connection, prefetchCount int, prefetchGlobal bool) (*rabbitMQChannel, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	rabbitCh := &rabbitMQChannel{
		uuid:       id.String(),
		connection: conn,
	}
	if err := rabbitCh.Connect(prefetchCount, prefetchGlobal); err != nil {
		return nil, err
	}
	return rabbitCh, nil

}

func (r *rabbitMQChannel) Connect(prefetchCount int, prefetchGlobal bool) error {
	var err error
	r.channel, err = r.connection.Channel()
	if err != nil {
		return err
	}
	err = r.channel.Qos(prefetchCount, 0, prefetchGlobal)
	if err != nil {
		return err
	}
	err = r.channel.Confirm(false)
	if err != nil {
		return err
	}
	r.confirms = r.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	return nil
}

func (r *rabbitMQChannel) Close() error {
	if r.channel == nil {
		return errors.New("channel is nil")
	}
	return r.channel.Close()
}

func (r *rabbitMQChannel) Publish(exchange, key string, message amqp.Publishing) error {
	if r.channel == nil {
		return errors.New("channel is nil")
	}
	if err := r.channel.Publish(exchange, key, false, false, message); err != nil {
		return err
	}
	if confirmed := <-r.confirms; confirmed.Ack {
		return nil
	} else {
		return fmt.Errorf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

func (r *rabbitMQChannel) DeclareExchange(exchange string) error {
	return r.channel.ExchangeDeclare(
		exchange, // name
		"topic",  // kind
		false,    // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
}

func (r *rabbitMQChannel) DeclareDurableExchange(exchange string) error {
	return r.channel.ExchangeDeclare(
		exchange, // name
		"topic",  // kind
		true,     // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
}

func (r *rabbitMQChannel) DeclareQueue(queue string, args amqp.Table) error {
	_, err := r.channel.QueueDeclare(
		queue, // name
		false, // durable
		true,  // autoDelete
		false, // exclusive
		false, // noWait
		args,  // args
	)
	return err
}

func (r *rabbitMQChannel) DeclareDurableQueue(queue string, args amqp.Table) error {
	_, err := r.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		args,  // args
	)
	return err
}

func (r *rabbitMQChannel) DeclareReplyQueue(queue string) error {
	_, err := r.channel.QueueDeclare(
		queue, // name
		false, // durable
		true,  // autoDelete
		true,  // exclusive
		false, // noWait
		nil,   // args
	)
	return err
}

func (r *rabbitMQChannel) ConsumeQueue(queue string, autoAck bool) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queue,   // queue
		r.uuid,  // consumer
		autoAck, // autoAck
		false,   // exclusive
		false,   // nonlocal
		false,   // nowait
		nil,     // args
	)
}

func (r *rabbitMQChannel) BindQueue(queue, key, exchange string, args amqp.Table) error {
	return r.channel.QueueBind(
		queue,    // name
		key,      // key
		exchange, // exchange
		false,    // noWait
		args,     // args
	)
}
