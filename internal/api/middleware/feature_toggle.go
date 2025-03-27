package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func FeatureToggleAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if os.Getenv("DISABLE_USER_INPUT") == "true" {
			renderForbidden(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
