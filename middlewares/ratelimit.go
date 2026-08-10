package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiterMiddleware 限制单个 IP 每分钟最多访问 limit 次
func RateLimiterMiddleware(rdb *redis.Client, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		// 构造限流 Key，按分钟隔离：rate:127.0.0.1:202608101004
		now := time.Now()
		rateKey := fmt.Sprintf("rate:%s:%s", ip, now.Format("200601021504"))

		reqCtx := c.Request.Context()

		// 1. 递增当前分钟的计数
		count, err := rdb.Incr(reqCtx, rateKey).Result()
		if err != nil {
			// 如果 Redis 出错，通常选择放行（降级），避免影响正常用户
			c.Next()
			return
		}

		// 2. 如果是当前分钟的第 1 次请求，设置 1 分钟自动过期，防止内存泄漏
		if count == 1 {
			rdb.Expire(reqCtx, rateKey, time.Minute)
		}

		// 3. 超过限制，直接终止请求
		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Limit is 3 requests per minute.",
			})
			return // 必须 return，阻止执行 c.Next()
		}

		// 4. 未超限，继续执行后续的主逻辑
		c.Next()
	}
}
