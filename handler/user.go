package handler

import (
	"net/http"
	"rekber/helper"
	"rekber/model"
	"rekber/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(services *service.UserService) *UserHandler {
	return &UserHandler{service: services}
}

func (u *UserHandler) Register(c *gin.Context) {
	req := model.User{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.ResponseAPI(false, http.StatusBadRequest, err.Error(), model.User{}))
		return
	}

	newUser, err := u.service.Register(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), model.User{}))
		return
	}

	c.JSON(http.StatusOK, helper.ResponseAPI(true, http.StatusOK, "success register new user", newUser))
	return
}
