package main

import (
	"os"

	"github.com/BieggerM/userservice/pkg/adapter/out/broker"
	"github.com/BieggerM/userservice/pkg/adapter/out/database"
	"github.com/BieggerM/userservice/pkg/models"
	"github.com/sirupsen/logrus"

)

var DB database.Postgres
var MB broker.RabbitMQ

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Check Connection to Message Broker
	prepareBroker()

	// Connect to PostgreSQL
	// Run migrations
	prepareDatabase()

	// Create demo users
	createDemoUsers()

	// Start the Gin server
	startGinServer()

	// Start GRPC
	StartGRPCServer()

}

func prepareBroker() {
	if err := MB.Connect(
		os.Getenv("RABBIT_USER"),
		os.Getenv("RABBIT_PASSWORD"),
		os.Getenv("RABBIT_HOST"),
		os.Getenv("RABBIT_PORT")); err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	} else {
		logrus.Infoln("Connected to RabbitMQ")
	}
	defer MB.Disconnect()
}

func prepareDatabase() {
	if dberr := DB.Connect(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME")); dberr != nil {
		logrus.Fatalf("Failed to connect to PostgreSQL: %v", dberr)
	} else {
		logrus.Info("Connected to PostgreSQL")
	}
	defer DB.Close()

	if err := DB.RunMigrations("file://migrations"); err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	} else {
		logrus.Info("Migrations run successfully")
	}
}



func createDemoUsers() {
	demoUsers := []models.User{
		{Username: "user1", FirstName: "John", LastName: "Doe"},
		{Username: "user2", FirstName: "Jane", LastName: "Doe"},
		{Username: "user3", FirstName: "Jim", LastName: "Beam"},
	}

	for _, user := range demoUsers {
		if err := DB.SaveUser(user); err != nil {
			logrus.Errorf("Failed to create demo user %s: %v", user.Username, err)
		} else {
			logrus.Infof("Created demo user %s", user.Username)
		}
	}
}
