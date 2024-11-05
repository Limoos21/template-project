package middleware

import (
	"github.com/gin-gonic/gin"
)

// NoOpMiddleware — это middleware, которое ничего не делает.
func NoOpMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Здесь можно добавить логику в будущем, если понадобится

		// Вызов следующего обработчика
		c.Next()
	}
}
