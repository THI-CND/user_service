package broker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
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
func (r *RabbitMQ) Connect() error {
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
    rabbitmqPort := os.Getenv("RABBITMQ_PORT")
    rabbitmqUser := os.Getenv("RABBITMQ_USER")
    rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")
	var err error
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort)
	r.Conn, err = amqp.Dial(connStr)
    if err != nil {
        fmt.Println("Failed to connect to RabbitMQ:", err)
        return err
    }
    fmt.Println("Connected to RabbitMQ")
    return nil
}

// Disconnect disconnects from the RabbitMQ message broker
func (r *RabbitMQ) Disconnect() error {
	return r.Conn.Close()
}

// Publish publishes a message to the RabbitMQ message broker
func (r *RabbitMQ) Publish(exchange, key string, msgBody []byte) error {
	message := amqp.Publishing{
        ContentType: "application/json", // Set content type for clarity
        Body:        msgBody,
    }
	
	ch, err := r.Conn.Channel()
	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)

	q, err := ch.QueueDeclare(
		"users", // Queue name
		true,        // Durable
		false,       // Auto-delete
		false,       // Exclusive
		false,       // No-wait
		nil,         // Arguments
	)
	q = q
	

	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.Publish(exchange, key, false, false, message) 
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