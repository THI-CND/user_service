package restserver

import (
	"encoding/json"

	"github.com/BieggerM/userservice/pkg/adapter/out/broker"
	"github.com/BieggerM/userservice/pkg/adapter/out/database"
	"github.com/BieggerM/userservice/pkg/adapter/out/logger"
	"github.com/BieggerM/userservice/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RestServer interface {
	StartRestServer(MB broker.MessageBroker, DB database.Database, rlog logger.Logger)
}
type GinServer struct {
	DB   database.Database
	MB   broker.MessageBroker
	rlog logger.Logger
}

func (g *GinServer) StartRestServer(MB broker.MessageBroker, DB database.Database, rlog logger.Logger) {
	g.DB = DB
	g.MB = MB
	g.rlog = rlog
	r := gin.Default()
	userGroup := r.Group("/api/v1/users")
	userGroup.GET("", g.listUsers)
	userGroup.GET("/:username", g.getUser)
	userGroup.POST("", g.createUser)
	userGroup.PATCH("", g.updateUser)
	userGroup.DELETE("", g.deleteUser)
	logrus.Infof("Gin Server started on port %s", ":8082")
	if err := r.Run(":8082"); err != nil {
		logrus.Fatalf("Failed to run Gin server: %v", err)
	}
}

func (g *GinServer) listUsers(c *gin.Context) {
	users := g.DB.ListUsers()
	c.JSON(200, gin.H{
		"users": users,
	})
}

func (g *GinServer) getUser(c *gin.Context) {
	user, err := g.DB.GetUser(c.Param("username"))
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	c.JSON(200, gin.H{
		"username":  user.Username,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
	})
}

func (g *GinServer) createUser(c *gin.Context) {
	var user models.User
	c.ShouldBindBodyWithJSON(&user)
	if err := g.DB.SaveUser(user); err != nil {
		c.JSON(500, gin.H{"error": "failed to save user to database - username exists"})
		return
	}
	if err := g.publishEvents(user, c); err != nil {
		c.JSON(500, gin.H{"error": "failed to publish events to RabbitMQ"})
		return
	}
	g.rlog.Info("User created", "username", user.Username)
}

func (g *GinServer) updateUser(c *gin.Context) {
	var user models.User
	c.ShouldBindBodyWithJSON(&user)
	oldUser, err := g.DB.GetUser(user.Username)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	g.DB.UpdateUser(user)
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
	if err := g.MB.Publish("recipemanagement", "users.update", msgBody); err != nil {
		c.JSON(500, gin.H{"error": "failed to publish message to RabbitMQ"})
	}
}

func (g *GinServer) deleteUser(c *gin.Context) {
	var user models.User
	c.ShouldBindBodyWithJSON(&user)
	g.DB.DeleteUser(user.Username)
	c.JSON(200, gin.H{
		"message":  "user deleted",
		"username": user.Username,
	})
	g.rlog.Info("User deleted", "username", user.Username)
}

func (g *GinServer) publishEvents(user models.User, c *gin.Context) error {
	// Prepare message for RabbitMQ
	// Marshall user struct to JSON
	// publish message to RabbitMQ exchange user with routing key "users.new"
	// publish message to RabbitMQ exchange user with routing key "users.count"
	msgBody, err := json.Marshal(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal user to JSON"})
	}

	if err := g.MB.Publish("recipemanagement", "users.new", msgBody); err != nil {
		c.JSON(500, gin.H{"error": "failed to publish message to RabbitMQ"})
	}

	userCount := len(g.DB.ListUsers())
	userCountBytes, err := json.Marshal(userCount)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal user count to JSON"})
	}
	if err := g.MB.Publish("recipemanagement", "users.count", userCountBytes); err != nil {
		c.JSON(500, gin.H{"error": "failed to publish message to RabbitMQ"})
	}

	return nil
}
