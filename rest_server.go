package main

import (
    "encoding/json"
    "github.com/BieggerM/userservice/pkg/models"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

func startGinServer() {
    r := gin.Default()
    user_group := r.Group("/api/v1/users")
    user_group.GET("", listUsers)
    user_group.GET("/:username", getUser)
    user_group.POST("", createUser)
    user_group.PATCH("", updateUser)
    user_group.DELETE("", deleteUser)
    if err := r.Run(":8082"); err != nil {
        logrus.Fatalf("Failed to run Gin server: %v", err)
    }
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
        "username":  user.Username,
        "firstname": user.FirstName,
        "lastname":  user.LastName,
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

func updateUser(c *gin.Context) {
    var user models.User
    c.ShouldBindBodyWithJSON(&user)
    oldUser := DB.GetUser(user.Username)

    DB.UpdateUser(user)
    c.JSON(200, gin.H{
        "message":   "user updated",
        "username":  user.Username,
        "firstname": user.FirstName,
        "lastname":  user.LastName,
    })

    msgBody, err := json.Marshal(map[string]interface{}{
        "oldUser":     oldUser,
        "updatedUser": user,
    })
    if err != nil {
        c.JSON(500, gin.H{"error": "failed to marshal user to JSON"})
    }
    if err := MB.Publish("recipemanagement", "users.update", msgBody); err != nil {
        c.JSON(500, gin.H{"error": "failed to publish message to RabbitMQ"})
    }
}

func deleteUser(c *gin.Context) {
    var user models.User
    c.ShouldBindBodyWithJSON(&user)
    DB.DeleteUser(user.Username)
    c.JSON(200, gin.H{
        "message":  "user deleted",
        "username": user.Username,
    })
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
