package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Joke will contain information about a joke
type Joke struct {
	ID      int    `json:"id" binding:"required"`
	Likes   int    `json:"likes"`
	Content string `json:"content" binding:"required"`
}

var jokes = []Joke{
	Joke{1, 0, "JOKE 1"},
	Joke{2, 0, "JOKE 2"},
	Joke{3, 0, "JOKE 3"},
	Joke{4, 0, "JOKE 4"},
	Joke{5, 0, "JOKE 5"},
	Joke{6, 0, "JOKE 6"},
	Joke{7, 0, "JOKE 7"},
	Joke{8, 0, "JOKE 8"},
}

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
	c.JSON(http.StatusOK, jokes)
}

// LikeJoke increments the likes of a specific Joke
// and will return the Joke struct or a Status Not Found
func LikeJoke(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	if jokeId, err := strconv.Atoi(c.Param("jokeID")); err == nil {
		for _, joke := range jokes {
			if joke.ID == jokeId {
				joke.Likes++
				c.JSON(http.StatusOK, &joke)
				break
			}
		}
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}
