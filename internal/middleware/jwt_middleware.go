package middleware

import (
	"net/http" 
	"strings"

	"github.com/artembliss/go-fitness-tracker/internal/services"
	"github.com/artembliss/go-fitness-tracker/pkg/auth"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(s *services.UserService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == ""{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			ctx.Abort()
			return
		}
		tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
            ctx.Abort()
            return
        }
		token := tokenParts[1]
        claims, err := auth.VerifyJWT(token)
        if err != nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            ctx.Abort()
            return
        }

		email := claims.Subject

        user, err := s.GetUserByEmail(email)
        if err != nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            ctx.Abort()
            return
        }
		ctx.Set("userID", user.ID)

        ctx.Next()
	}

}