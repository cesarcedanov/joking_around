package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Joke will contain information about a joke
type Joke struct {
	ID      int    `json:"id" binding:"required"`
	Likes   int    `json:"likes"`
	Content string `json:"content" binding:"required"`
}

var jokes = []*Joke{
	&Joke{1, 0, "JOKE 1"},
	&Joke{2, 0, "JOKE 2"},
	&Joke{3, 0, "JOKE 3"},
	&Joke{4, 0, "JOKE 4"},
	&Joke{5, 0, "JOKE 5"},
	&Joke{6, 0, "JOKE 6"},
	&Joke{7, 0, "JOKE 7"},
	&Joke{8, 0, "JOKE 8"},
}

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

var jwtMiddleWare *jwtmiddleware.JWTMiddleware

func main() {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: validationKey,
		SigningMethod:       jwt.SigningMethodRS256,
	})

	// register our actual jwtMiddleware
	jwtMiddleWare = jwtMiddleware

	// Create a router from the Gin's Default router
	router := gin.Default()

	// Serve the frontend static files
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// Create a group of router with the path prefix '/api/'
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/", HealthCheck)

		apiRoutes.GET("/jokes", authMiddleware(), GetJokes)
		apiRoutes.POST("/jokes/like/:jokeID", authMiddleware(), LikeJoke)

	}

	// Start and Run the Server on port 9000
	router.Run(":9000")
}

func validationKey(token *jwt.Token) (interface{}, error) {
	// verify the api audience
	aud := os.Getenv("AUTH0_API_AUDIENCE")
	checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
	if !checkAudience {
		return token, errors.New("Invalid audience.")
	}
	// verify iss claim
	iss := os.Getenv("AUTH0_DOMAIN")
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	if !checkIss {
		return token, errors.New("Invalid issuer.")
	}
	cert, err := getPemCert(token)
	if err != nil {
		log.Fatalf("could not get cert: %+v", err)
	}
	result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	return result, nil
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(os.Getenv("AUTH0_DOMAIN") + ".well-known/jwks.json")
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()
	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}
	x5c := jwks.Keys[0].X5c
	for k, v := range x5c {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + v + "\n-----END CERTIFICATE-----"
		}
	}
	if cert == "" {
		return cert, errors.New("unable to find appropriate key.")
	}
	return cert, nil
}

// authMiddleware intercepts the requests, and check for a valid jwt token
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the client secret key
		err := jwtMiddleWare.CheckJWT(c.Writer, c.Request)
		if err != nil {
			// Token not found
			fmt.Println(err)
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
			return
		}
	}
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
