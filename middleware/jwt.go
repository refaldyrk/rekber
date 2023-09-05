package middleware

import (
	"net/http"
	"rekber/helper"
	"rekber/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"
)

func JWTMiddleware(db *qmgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header.Get("Authorization")

		tokenCookies, _ := c.Cookie("Authorization")

		if authorizationHeader == "" && tokenCookies == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
				"error":   true,
			})
			c.Abort()
			return
		}

		tokenStringCookie := tokenCookies
		tokenStringHeader := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		claims, err := helper.ValidateJWT(tokenStringCookie)
		if err != nil {
			claims, err = helper.ValidateJWT(tokenStringHeader)
			if err != nil {

				c.JSON(http.StatusUnauthorized, gin.H{
					"message": err.Error(),
					"error":   true,
				})
				c.Abort()
				return
			}
		}

		userID := claims["sub"].(string)
		user := model.User{}
		err = db.Collection("User").Find(c, qmgo.M{"user_id": userID}).One(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized, user not found",
				"error":   true,
			})
			c.Abort()
			return
		}

		c.Set("userID", user.UserID)
		c.Set("user", user)
		c.Next()
	}
}
