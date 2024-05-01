package middleware

import (
	"basicapi/config"
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

func APIKeyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader("X-API-KEY") != config.GetConfig().ApiKey {
			ctx.JSON(401, gin.H{"error": "API key invalid"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-KEY")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}

}

func LoggerMiddleware() gin.HandlerFunc {
	// TODO in case of a real project we should use a jwt of some kind to be able to identify the user and log the request with the user id.
	return func(ctx *gin.Context) {
		ctx.Next()
		config.Logger.Info().Msgf("Request received: %s %s at addr: %s", ctx.Request.Method, ctx.Request.URL, ctx.Request.RemoteAddr)
	}
}

var visitors = make(map[string]*rate.Limiter)
var mtx sync.Mutex

// Get a rate limiter for a given IP address.

func RateLimitMiddleware(c *gin.Context) {
	ip := c.ClientIP()

	mtx.Lock()
	defer mtx.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// Allow 1 request per second with a burst of 5 requests.
		limiter = rate.NewLimiter(1, 5)
		visitors[ip] = limiter
	}

	if !limiter.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": "too many requests",
		})
		return
	}

	c.Next()
}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel() // Important to avoid a context leak

		// Replace the request with a context-aware one
		c.Request = c.Request.WithContext(ctx)

		// Set up a channel to catch the completion signal
		finish := make(chan struct{})
		go func() {
			c.Next() // Process request
			finish <- struct{}{}
		}()

		// Monitor for the finish signal or the timeout
		select {
		case <-finish:
			// If finished, just return
			return
		case <-ctx.Done():
			// If the context is done, it must have timed out
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"error": "request timed out"})
		}
	}
}
