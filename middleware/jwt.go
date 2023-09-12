package middleware

import (
	"net/http"
	"rekber/helper"
	"rekber/model"
	"rekber/repository"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"
)

func JWTMiddleware(db *qmgo.Database, authRepo *repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Check Token Is Logout
		authorizationHeader := c.Request.Header.Get("Authorization")

		logout, err := authRepo.FindLogoutByToken(c, authorizationHeader)
		if !logout.ID.IsZero() {
			c.JSON(http.StatusUnauthorized, helper.ResponseAPI(false, http.StatusUnauthorized, "token has logout", gin.H{
				"time": logout.LogoutAt,
			}))
			c.Abort()
			return
		}

		if err != nil && err != qmgo.ErrNoSuchDocuments {
			c.JSON(http.StatusUnauthorized, helper.ResponseAPI(false, http.StatusUnauthorized, err.Error(), gin.H{}))
			c.Abort()
			return
		}

		tokenCookies, _ := c.Cookie("Authorization")

		if authorizationHeader == "" && tokenCookies == "" {
			c.JSON(http.StatusUnauthorized, helper.ResponseAPI(false, http.StatusUnauthorized, "unauthorized", gin.H{}))
			c.Abort()
			return
		}

		tokenStringCookie := tokenCookies
		tokenStringHeader := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		claims, err := helper.ValidateJWT(tokenStringCookie)
		if err != nil {
			claims, err = helper.ValidateJWT(tokenStringHeader)
			if err != nil {

				c.JSON(http.StatusUnauthorized, helper.ResponseAPI(false, http.StatusUnauthorized, err.Error(), gin.H{}))
				c.Abort()
				return
			}
		}

		userID := claims["sub"].(string)
		user := model.User{}
		err = db.Collection("User").Find(c, qmgo.M{"user_id": userID}).One(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, helper.ResponseAPI(false, http.StatusUnauthorized, "user not found", gin.H{}))
			c.Abort()
			return
		}

		c.Set("userID", user.UserID)
		c.Set("user", user)
		c.Next()
	}
}
