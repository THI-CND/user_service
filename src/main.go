package main

import (
	"github.com/BieggerM/userservice/pkg/adapter/in/grpcserver"
	"github.com/BieggerM/userservice/pkg/adapter/in/restserver"
	"github.com/BieggerM/userservice/pkg/adapter/out/authenticationprovider"
	"github.com/BieggerM/userservice/pkg/adapter/out/broker"
	"github.com/BieggerM/userservice/pkg/adapter/out/database"
	"github.com/BieggerM/userservice/pkg/adapter/out/logger"
	"github.com/BieggerM/userservice/pkg/models"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// Levels defines which log levels trigger the hook.

var DB database.Database
var MB broker.MessageBroker
var RS restserver.RestServer
var GS grpcserver.GrpcServer
var rlog logger.Logger
var auth authenticationprovider.AuthenticationProvider

func main() {
	// Initialize Implementations
	rlog = &logger.RemoteLogger{}
	DB = &database.Postgres{}
	MB = &broker.RabbitMQ{}
	RS = &restserver.GinServer{}
	GS = &grpcserver.UserServiceServer{}
	auth = &authenticationprovider.AuthProvider{}

	// Setup remote logging
	setupRemoteLogging()
	defer rlog.Close()

	// Check Connection to Message Broker
	prepareBroker()
	defer MB.Close()

	// Connect to PostgreSQL
	// Run migrations
	prepareDatabase()
	defer DB.Close()

	// Create demo users
	createDemoUsers()

	// Start the Gin server
	go RS.StartRestServer(MB, DB, rlog, auth)

	// Start GRPC
	GS.StartGRPCServer(MB, DB, rlog)

}

func setupRemoteLogging() {
	fluentPort, ferr := strconv.Atoi(os.Getenv("FLUENTD_PORT"))
	if ferr != nil {
		logrus.Fatalf("Invalid fluentd port number: %v", ferr)
	}
	var err error
	err = rlog.Setup(os.Getenv("FLUENTD_HOST"), fluentPort, "user-service")
	if err != nil {
		logrus.Fatalf("Failed to create remote logger: %v", err)
	}
}

func prepareBroker() {
	if err := MB.Connect(
		os.Getenv("RABBIT_USER"),
		os.Getenv("RABBIT_PASSWORD"),
		os.Getenv("RABBIT_HOST"),
		os.Getenv("RABBIT_PORT")); err != nil {
		rlog.Fatal("Failed to connect to MessageBroker:", "error", err)
	} else {
		rlog.Info("Connected to MessageBroker", "host", os.Getenv("RABBIT_HOST"), "port", os.Getenv("RABBIT_PORT"))
	}
}

func prepareDatabase() {
	if dberr := DB.Connect(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME")); dberr != nil {
		rlog.Fatal("Failed to connect to Database: ", "error", dberr)
	} else {
		rlog.Info("Connected to Database", "host", os.Getenv("DB_HOST"), "port", os.Getenv("DB_PORT"))
	}

	if err := DB.RunMigrations("file://migrations"); err != nil {
		rlog.Warn("Failed to run migrations", "error", err)
	} else {
		rlog.Info("Migrations run successfully", "path", "file://migrations")
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
			logrus.Warn("Failed to create demo user", "", "")
		} else {
			logrus.Infof("Created demo user %s", user.Username)
		}
	}
}
