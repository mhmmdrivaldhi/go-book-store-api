package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/service"
)

type AuthMiddleware interface {
	RequireToken() gin.HandlerFunc
}

type authMiddleware struct{
	jwtService service.JwtService
}

func (am *authMiddleware) RequireToken() gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		log.Println("Middleware is running...")

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization is missing header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Authorization header is requiresd"})
			return 
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			log.Println("Bearer token is missing")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Bearer token is required"})
			return
		}

		claims, err := am.jwtService.ValidateToken(tokenStr)
		if err != nil {
			log.Printf("Middleware: Token validation/parsing failed: %v\n", err)

			if errors.Is(err, errors.New("token is expire")) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Token has expired"})
			} else {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid token"})
			}
			return 
		}

		if claims.UserID == 0 {
			log.Printf("Middleware: Warning - Parsed UserID is 0. Claims: %+v\n", claims)
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{
		jwtService: jwtService,
	}
}