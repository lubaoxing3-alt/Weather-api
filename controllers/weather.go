package controllers

import (
	"net/http"
	"weather-api/services"

	"github.com/gin-gonic/gin"
)

type WeatherController struct {
	WeatherService *services.WeatherService
}

func NewWeatherController(ws *services.WeatherService) *WeatherController {
	return &WeatherController{WeatherService: ws}
}

func (ctrl *WeatherController) GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
		return
	}

	data, statusCode, err := ctrl.WeatherService.GetWeatherData(c.Request.Context(), city)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.Data(statusCode, "application/json", data)
}
