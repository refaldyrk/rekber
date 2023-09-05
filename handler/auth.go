package handler

import (
	"net/http"
	"rekber/dto"
	"rekber/helper"
	"rekber/service"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(services *service.AuthService) *AuthHandler {
	return &AuthHandler{service: services}
}

func (u *AuthHandler) Login(ctx *gin.Context) {
	var loginReq dto.LoginReq
	err := ctx.ShouldBindJSON(&loginReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ResponseAPI(false, http.StatusBadRequest, err.Error(), gin.H{}))
		return
	}

	user, err := u.service.Login(ctx, loginReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

	// save to db

	// generate token
	token, err := helper.GenJWT(user.UserID, 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, helper.ResponseAPI(true, http.StatusOK, "success login", gin.H{
		"token": token,
		"user":  user,
	}))
	return
}
