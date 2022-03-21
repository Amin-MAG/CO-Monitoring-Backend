package middleware

import "github.com/gin-gonic/gin"

type GinMiddleware interface {
	Middleware() gin.HandlerFunc
}
