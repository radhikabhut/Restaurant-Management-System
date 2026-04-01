package middleware

import (
	"net/http"
	"restaurant-management-system/handler"
	"restaurant-management-system/objects"
	"restaurant-management-system/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, objects.GenericError{
				Success:    false,
				ErrMsg:     "authorization header is required",
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, objects.GenericError{
				Success:    false,
				ErrMsg:     "authorization header format must be Bearer <token>",
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, objects.GenericError{
				Success:    false,
				ErrMsg:     "invalid or expired token",
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func RBACMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, objects.GenericError{
				Success:    false,
				ErrMsg:     "user role not found in context",
				StatusCode: http.StatusForbidden,
			})
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, objects.GenericError{
				Success:    false,
				ErrMsg:     "user role is not a valid string",
				StatusCode: http.StatusForbidden,
			})
			return
		}

		isAllowed, err := handler.CheckPermission(c.Request.Context(), roleStr, permission)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, objects.GenericError{
				Success:    false,
				ErrMsg:     "failed to check permissions",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, objects.GenericError{
				Success:    false,
				ErrMsg:     "you do not have permission to access this resource",
				StatusCode: http.StatusForbidden,
			})
			return
		}

		c.Next()
	}
}

