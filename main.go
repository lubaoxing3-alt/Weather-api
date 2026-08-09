package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	r := gin.Default()

	r.GET("/weather", func(c *gin.Context) {
		_ = godotenv.Load()
		location := url.PathEscape(c.Query("city"))
		if location == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
			return
		}
		apiKey := os.Getenv("VISUAL_CROSSING_KEY")
		reqURL := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s", location, apiKey)

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
		c.Data(http.StatusOK, "application/json", body)
	})
	r.Run()
}
