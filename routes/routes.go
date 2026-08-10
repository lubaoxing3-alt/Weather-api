package routes

import (
	"weather-api/controllers"
	"weather-api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(rdb *redis.Client, weatherCtrl *controllers.WeatherController) *gin.Engine {
	r := gin.Default()

	// 挂载限流中间件
	r.Use(middlewares.RateLimiterMiddleware(rdb, 3))

	r.GET("/weather", weatherCtrl.GetWeather)

	return r
}
