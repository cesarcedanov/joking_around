package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a router from the Gin's Default router
	router := gin.Default()

	// Create a group of router with the path prefix '/api/'
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello World",
			})
		})
	}

	// Start and Run the Server on port 9000
	router.Run(":9000")
}
