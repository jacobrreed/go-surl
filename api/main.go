package main

import (
	"errors"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Url struct {
	Url string `json:"url"`
}

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = uint64(len(alphabet))
)

// function to recieve GET request that accepts a url
func generateUrl(c *gin.Context) {
	// Retrieve originalUrl from request body json
	var originalUrl Url
	if err := c.ShouldBindJSON(&originalUrl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return a JSON response with the shortened url
	c.JSON(
		http.StatusOK,
		Url{Url: "https://neonvoid.io/g9Azxl3"})
}

func main() {
	// Setup Redis
	rdb := setupRedis()

	// Setup api router
	r := gin.Default()
	r.POST("/", generateUrl)
	r.Run()
}

func setupRedis() (client *redis.Client) {
	// setup redis
	redisHost := getenv("REDIS_HOST", "redis")
	redisPort := getenv("REDIS_PORT", "6379")
	redisPassword := getenv("REDIS_PASSWORD", "")
	redisDb := getenv("REDIS_DB", "0")

	redisDbNumber, err := strconv.Atoi(redisDb)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       redisDbNumber,
	})
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func Encode(number uint64) string {
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(11)

	for ; number > 0; number = number / length {
		encodedBuilder.WriteByte(alphabet[(number % length)])
	}

	return encodedBuilder.String()
}

func Decode(encoded string) (uint64, error) {
	var number uint64

	for i, symbol := range encoded {
		alphabeticPosition := strings.IndexRune(alphabet, symbol)

		if alphabeticPosition == -1 {
			return uint64(alphabeticPosition), errors.New("invalid character: " + string(symbol))
		}
		number += uint64(alphabeticPosition) * uint64(math.Pow(float64(length), float64(i)))
	}

	return number, nil
}
