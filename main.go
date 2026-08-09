package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/weather", func(c *gin.Context) {
		city := c.Query("city")
		if city == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
			return
		}

		// 第一步先返回 Mock 数据，验证 HTTP 通路
		c.JSON(http.StatusOK, gin.H{
			"city":        city,
			"temperature": "28°C",
			"status":      "Mock Data (Server OK)",
		})
	})
	r.Run()
}
