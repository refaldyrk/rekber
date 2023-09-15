package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rekber/constant"
	"rekber/helper"
)

func IsAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		role := context.GetString("role")
		if role == "" || role == constant.MEMBER_ROLE {
			context.AbortWithStatusJSON(http.StatusForbidden, helper.ResponseAPI(false, http.StatusForbidden, "access denied", gin.H{}))
			return
		}

		context.Next()
	}
}
