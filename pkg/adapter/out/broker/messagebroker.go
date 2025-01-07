package broker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageBroker is the interface for the message broker
type MessageBroker interface {
	Connect() error
	Disconnect() error
	Publish(exchange, key string, msg amqp.Publishing) error
	Subscribe(exchange, key string, handler func(amqp.Delivery)) error
}

// RabbitMQ is the RabbitMQ message broker
type RabbitMQ struct {
	Conn *amqp.Connection
}
// Connect connects to the RabbitMQ message broker
func (r *RabbitMQ) Connect(rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort string) error {
    connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort)
    var err error
    r.Conn, err = amqp.Dial(connStr)
    if err != nil {
        fmt.Println("Failed to connect to RabbitMQ:", err)
        return err
    }
    return nil
}

// Disconnect disconnects from the RabbitMQ message broker
func (r *RabbitMQ) Disconnect() error {
    if r.Conn != nil {
        return r.Conn.Close()
    }
    return nil
}

// Publish publishes a message to the RabbitMQ message broker
func (r *RabbitMQ) Publish(exchange, routingKey string, body []byte) error {
    if r.Conn == nil {
        return fmt.Errorf("connection is nil")
    }

    ch, err := r.Conn.Channel()
    if err != nil {
        return fmt.Errorf("failed to open a channel: %w", err)
    }
    defer ch.Close()

    if ch == nil {
        return fmt.Errorf("channel is nil")
    }

    err = ch.ExchangeDeclare(
        exchange, // name
        "topic",  // type
        true,     // durable
        false,    // auto-deleted
        false,    // internal
        false,    // no-wait
        nil,      // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare exchange: %w", err)
    }

    err = ch.Publish(
        exchange,   // exchange
        routingKey, // routing key
        false,      // mandatory
        false,      // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        })
    if err != nil {
        return fmt.Errorf("failed to publish message: %w", err)
    }

	return nil
}
// Subscribe subscribes to a message from the RabbitMQ message broker
// Implement GoRoutine to handle messages
func (r *RabbitMQ) Subscribe(exchange, key string) error {
	ch, err := r.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	msgs, err := ch.Consume(
		"users",
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

	for msg := range msgs {
		handlemsg(msg)
	}
	
	return nil
}

func handlemsg(msg amqp.Delivery) {
	// Do something with the message
	fmt.Println(string(msg.Body))
}