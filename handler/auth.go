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

func (u *AuthHandler) LoginV2Register(ctx *gin.Context) {
	var loginReq dto.LoginV2Req
	err := ctx.ShouldBindJSON(&loginReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ResponseAPI(false, http.StatusBadRequest, err.Error(), gin.H{}))
		return
	}

	loginV2, err := u.service.RegisterLoginV2(ctx, loginReq.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.ResponseAPI(true, http.StatusOK, "success login", loginV2))
	return
}

func (u *AuthHandler) LoginV2(ctx *gin.Context) {
	param := ctx.Param("codelink")
	if param == "" {
		ctx.JSON(http.StatusBadRequest, helper.ResponseAPI(false, http.StatusBadRequest, "codelink has not found", gin.H{}))
		return
	}

	user, err := u.service.LoginV2(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

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

func (u *AuthHandler) Logout(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		ctx.JSON(http.StatusUnprocessableEntity, helper.ResponseAPI(false, http.StatusUnprocessableEntity, "unauthorized", gin.H{}))
		return
	}

	err := u.service.Logout(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.ResponseAPI(true, http.StatusOK, "success logout", gin.H{
		"message": "success logout",
	}))
}
