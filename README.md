# User Service

This is a User Service built with Go, Gin, PostgreSQL, and RabbitMQ. It provides RESTful APIs to manage users, including creating, updating, and deleting users. The service also publishes user creation messages to RabbitMQ.

## Prerequisites

- Go 1.16+
- Docker
- Docker Compose

## Getting Started

### Running with Docker Compose

To run the service with Docker Compose, use the following command:

```sh
docker-compose up
```
This will start the PostgreSQL, RabbitMQ, and User Service containers.

## Running Locally
To run the service locally, set the following environment variables and run the application:

```sh
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=recipe
export RABBITMQ_HOST=localhost
export RABBITMQ_PORT=5672
export RABBITMQ_USER=guest
export RABBITMQ_PASSWORD=guest

go run main.go
```
## REST API
### Create User
URL: /users
Method: POST
Request Body:
 ```json
{
  "username": "johndoe",
  "firstname": "John",
  "lastname": "Doe"
}
```

### Update User
URL: /users
Method: PUT
Request Body:

```json
{
  "username": "johndoe",
  "firstname": "John",
  "lastname": "Doe"
}
```
### Delete User
URL: /users/:username
Method: DELETE
Success Response:
* Code: 200 OK

### Get User
URL: /users/:username
Method: GET
Success Response:
* Code: 200 OK

## MessageBroker
The User Service uses RabbitMQ as a message broker to publish user creation messages. A simple interface is provided to allow for similar tools. The RabbitMQ struct in the messagebroker.go file handles the connection, publishing, and subscribing to RabbitMQ.

### Connect
The Connect method establishes a connection to the RabbitMQ server using the provided credentials and host information.

### Disconnect
The Disconnect method closes the connection to the RabbitMQ server.

### Publish
The Publish method publishes a message to the specified exchange with the given routing key and message body.

### Subscribe
The Subscribe method subscribes to messages from the specified exchange and routing key, and processes them using the provided handler function.

Example usage in main.go:

