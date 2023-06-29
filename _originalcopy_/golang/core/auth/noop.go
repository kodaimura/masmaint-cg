package auth

import (
	"github.com/gin-gonic/gin"
)

func NoopAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {}
}