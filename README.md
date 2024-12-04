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

