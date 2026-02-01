package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserKey = "auth.user"

type AuthedUser struct {
	UserID string
}

func Middleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := ParseAccessToken(jwtSecret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			return
		}
		c.Set(CtxUserKey, AuthedUser{UserID: claims.UserID})
		c.Next()
	}
}

func MustGetUser(c *gin.Context) AuthedUser {
	v, ok := c.Get(CtxUserKey)
	if !ok {
		panic("missing auth user in context")
	}
	return v.(AuthedUser)
}
