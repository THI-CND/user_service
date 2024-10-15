/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
func (r *RabbitMQ) Connect() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	r.Conn = conn
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
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.Publish(exchange, key, false, false, message) 
}

// Subscribe subscribes to a message from the RabbitMQ message broker
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