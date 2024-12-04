package main

import (
	"github.com/gin-gonic/gin"
	"github.com/BieggerM/userservice/pkg/database"
	"github.com/BieggerM/userservice/pkg/models"
	"github.com/BieggerM/userservice/pkg/broker"
	"encoding/json"
	"os"
)

var DB database.Postgres
var MB broker.RabbitMQ


func main() {
    // os.Setenv("DB_HOST", "localhost")
    // os.Setenv("DB_PORT", "5432")
    // os.Setenv("DB_USER", "postgres")
    // os.Setenv("DB_PASSWORD", "password")
    // os.Setenv("DB_NAME", "recipe")
    // os.Setenv("RABBIT_HOST", "localhost")
    // os.Setenv("RABBIT_PORT", "5672")
    // os.Setenv("RABBIT_USER", "guest")
    // os.Setenv("RABBIT_PASSWORD", "guest")

    err := MB.Connect(
        os.Getenv("RABBIT_USER"),
        os.Getenv("RABBIT_PASSWORD"),
        os.Getenv("RABBIT_HOST"),
        os.Getenv("RABBIT_PORT"),
    )
    if err != nil {
        panic("Failed to connect to RabbitMQ")
    }
    defer MB.Disconnect()

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

	// Prepare message for RabbitMQ
	msgBody, err := json.Marshal(user) // Marshall user struct to JSON
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal user to JSON"})
		return
	}
	// publish message to RabbitMQ exchange user with routing key ""
	if err := MB.Publish("users", "", msgBody); err != nil {
        c.JSON(500, gin.H{"error": "failed to publish message to RabbitMQ"})
        return
    }
	c.JSON(200, gin.H{
		"message": "user created",
	})
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
