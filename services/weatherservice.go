package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type WeatherService struct {
	Rdb *redis.Client
}

func NewWeatherService(rdb *redis.Client) *WeatherService {
	return &WeatherService{Rdb: rdb}
}

func (s *WeatherService) GetWeatherData(ctx context.Context, city string) ([]byte, int, error) {
	city = strings.ToLower(city)
	cacheKey := "weather:" + city

	// 1. 查缓存
	val, err := s.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		return []byte(val), http.StatusOK, nil
	}

	// 2. 调第三方 API
	apiKey := os.Getenv("VISUAL_CROSSING_KEY")
	reqURL := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s", url.PathEscape(city), apiKey)

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to reach weather service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("invalid city or third-party API error")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// 3. 写入缓存
	s.Rdb.Set(ctx, cacheKey, string(body), 12*time.Hour)

	return body, http.StatusOK, nil
}
