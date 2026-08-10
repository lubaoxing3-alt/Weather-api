package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"weather-api/middlewares"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // 密码
		DB:       0,                // 数据库索引
	})
	defer rdb.Close()

	_ = godotenv.Load()

	r := gin.Default()
	r.Use(middlewares.RateLimiterMiddleware(rdb, 3))

	r.GET("/weather", func(c *gin.Context) {
		city := c.Query("city")
		if city == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
			return
		}
		city = strings.ToLower(city)
		cacheKey := "weather:" + city
		val, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			c.Data(http.StatusOK, "application/json", []byte(val))
			return
		}

		apiKey := os.Getenv("VISUAL_CROSSING_KEY")
		reqURL := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s", url.PathEscape(city), apiKey)

		resp, err := http.Get(reqURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reach weather service"})
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Invalid city or third-party API error"})
			return
		}
		body, _ := io.ReadAll(resp.Body)
		rdb.Set(ctx, cacheKey, string(body), 12*time.Hour)
		c.Data(http.StatusOK, "application/json", body)
	})
	r.Run()
}
