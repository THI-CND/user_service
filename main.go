package main

import (
	"github.com/gin-gonic/gin"
	"github.com/BieggerM/userservice/pkg/database"
	"github.com/BieggerM/userservice/pkg/models"
)

var DB database.Postgres

func main() {
	// create database implementation
	r := gin.Default()
	r.GET("/users", listUsers)
	r.GET("/users/:username", getUser)
	r.POST("/users", createUser)
	r.PATCH("/users", updateUser)
	r.DELETE("/users/", deleteUser)
	r.Run(":8080")
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
	DB.SaveUser(user)
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
