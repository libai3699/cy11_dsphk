package middleware

import (
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc:  allowDevOrigin,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Accept-Language"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func allowDevOrigin(origin string) bool {
	if allowConfiguredOrigin(origin) {
		return true
	}

	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}

	host := parsed.Hostname()
	if host == "localhost" || host == "127.0.0.1" {
		return true
	}
	if strings.HasPrefix(host, "10.") || strings.HasPrefix(host, "192.168.") {
		return true
	}
	if strings.HasPrefix(host, "172.") {
		parts := strings.Split(host, ".")
		if len(parts) > 1 {
			second, err := strconv.Atoi(parts[1])
			return err == nil && second >= 16 && second <= 31
		}
	}
	return false
}

func allowConfiguredOrigin(origin string) bool {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return false
	}

	defaultAllowed := []string{
		"https://dspadmin.wangwei.tech",
	}
	for _, allowed := range defaultAllowed {
		if origin == allowed {
			return true
		}
	}

	for _, allowed := range strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ",") {
		allowed = strings.TrimSpace(allowed)
		if allowed == "" {
			continue
		}
		if allowed == "*" || origin == allowed {
			return true
		}
		if strings.HasPrefix(allowed, "*.") {
			parsed, err := url.Parse(origin)
			if err != nil {
				continue
			}
			suffix := strings.TrimPrefix(allowed, "*")
			if strings.HasSuffix(parsed.Hostname(), suffix) {
				return true
			}
		}
	}
	return false
}
