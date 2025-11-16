package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		role, exists := ctx.Get("role")
		if !exists {
			ctx.AbortWithStatusJSON(403, dto.ErrorResponse{Error: "Forbidden"})
			return 
		}

		userRole := role.(string)
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(403, dto.ErrorResponse{Error: "access denied"})
	}
}