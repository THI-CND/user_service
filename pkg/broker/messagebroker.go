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

package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"

)

// MessageBroker is the interface for the message broker
type MessageBroker interface {
	Connect() error
	Disconnect() error
	Publish(ctx context.Context, exchange, key string, msg amqp.Publishing) error
	Subscribe(ctx context.Context, exchange, key string, handler func(amqp.Delivery)) error
}

// RabbitMQ is the RabbitMQ message broker
type RabbitMQ struct {
	Conn *amqp.Connection
}

