package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a router from the Gin's Default router
	router := gin.Default()

	// Serve the frontend static files
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// Create a group of router with the path prefix '/api/'
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/", HealthCheck)

		apiRoutes.GET("/jokes", GetJokes)
		apiRoutes.POST("/jokes/like/:jokeID", LikeJoke)

	}

	// Start and Run the Server on port 9000
	router.Run(":9000")

}

// HealthCheck ping the server and expect a StatusOK
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Healt Check",
	})
}

// GetJokes return a list of Jokes
func GetJokes(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": "This endpoint will show a list of jokes",
	})
}

// LikeJoke return the total of likes of a specific Joke
func LikeJoke(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": "This endpoint will give a like to a joke.",
	})
}
