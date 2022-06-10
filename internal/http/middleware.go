package http

import (
	"net/http"

	"github.com/AtCliffUnderline/golang-diploma/internal/entities"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ur entities.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("auth")
		if cookie == "" || err != nil {
			if err != nil {
				c.String(http.StatusUnauthorized, "")
				c.Abort()
			}
		} else {
			user, err := ur.GetByToken(cookie)
			if err != nil {
				c.String(http.StatusInternalServerError, "")
				c.Abort()
			}
			c.Set("user", user)
		}
		c.Next()
	}
}
