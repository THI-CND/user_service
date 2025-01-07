package main

import (
	"github.com/gin-gonic/gin"
	"github.com/BieggerM/userservice/pkg/adapter/out/database"
	"github.com/BieggerM/userservice/pkg/models"
	"github.com/BieggerM/userservice/pkg/adapter/out/broker"
	"encoding/json"
	"os"
	"github.com/sirupsen/logrus"
)

var DB database.Postgres
var MB broker.RabbitMQ


func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.SetOutput(os.Stdout)
    logrus.SetLevel(logrus.InfoLevel)
	
	// Connect to RabbitMQ
    if err := MB.Connect(
		os.Getenv("RABBIT_USER"), 
		os.Getenv("RABBIT_PASSWORD"), 
		os.Getenv("RABBIT_HOST"), 
		os.Getenv("RABBIT_PORT"), ); 
		err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	} else {
		logrus.Infoln("Connected to RabbitMQ")
	}
    defer MB.Disconnect()

	  // Connect to PostgreSQL
	if dberr := DB.Connect(
		os.Getenv("DB_HOST"), 
		os.Getenv("DB_PORT"), 
		os.Getenv("DB_USER"), 
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_NAME")); 
		dberr != nil {
		logrus.Fatalf("Failed to connect to PostgreSQL: %v", dberr)
	} else {
		logrus.Info("Connected to PostgreSQL")
	}
	defer DB.Close()

	// Run migrations
	if err := DB.RunMigrations("file://migrations"); err != nil {
        logrus.Fatalf("Failed to run migrations: %v", err)
    } else {
		logrus.Info("Migrations run successfully")
	}

	createDemoUsers()
    
	// Start the Gin server
    r := gin.Default()
    r.GET("/users", listUsers)
    r.GET("/users/:username", getUser)
    r.POST("/users", createUser)
    r.PATCH("/users", updateUser)
    r.DELETE("/users/", deleteUser)
    r.Run(":8082")
}


func listUsers(c *gin.Context) {
	users := DB.ListUsers()
	c.JSON(200, gin.H{
		"users": users,
	})
}

func getUser(c *gin.Context) {
	user := DB.GetUser(c.Param("username"))
	c.JSON(200, gin.H{
		"username" : user.Username,
		"firstname" : user.FirstName,
		"lastname" : user.LastName,
	})
}

func createUser(c *gin.Context) {
    var user models.User
    c.ShouldBindBodyWithJSON(&user)
    if err := DB.SaveUser(user); err != nil {
        c.JSON(500, gin.H{"error": "failed to save user to database - username exists"})
        return
    }
    if err := publishEvents(user, c); err != nil {
		c.JSON(500, gin.H{"error": "failed to publish events to RabbitMQ"})
		return
	}
}

func publishEvents(user models.User, c *gin.Context) error {
	// Prepare message for RabbitMQ
    // Marshall user struct to JSON
    // publish message to RabbitMQ exchange user with routing key "users.new"
    // publish message to RabbitMQ exchange user with routing key "users.count"
	msgBody, err := json.Marshal(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal user to JSON"})
	}

	if err := MB.Publish("recipemanagement", "users.new", msgBody); err != nil {
		c.JSON(500, gin.H{"error": "failed to publish message to RabbitMQ"})
	}

	userCount := len(DB.ListUsers())
	userCountBytes, err := json.Marshal(userCount)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal user count to JSON"})
	}
	if err := MB.Publish("recipemanagement", "users.count", userCountBytes); err != nil {
		c.JSON(500, gin.H{"error": "failed to publish message to RabbitMQ"})
	}

	return nil
}

func updateUser(c *gin.Context) {
	var user models.User
	c.ShouldBindBodyWithJSON(&user)
	DB.UpdateUser(user)
	c.JSON(200, gin.H{
		"message": "user updated",
		"username" : user.Username,
		"firstname" : user.FirstName,
		"lastname" : user.LastName,
	})
}

func deleteUser(c *gin.Context) {
	var user models.User
	c.ShouldBindBodyWithJSON(&user)
	DB.DeleteUser(user.Username)
	c.JSON(200, gin.H{
		"message": "user deleted",
		"username" : user.Username,
	})
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
