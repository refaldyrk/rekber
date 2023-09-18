package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rekber/helper"
	"rekber/service"
)

type BalanceHandler struct {
	service *service.BalanceService
}

func NewBalanceHandler(serv *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		service: serv,
	}
}

func (b *BalanceHandler) FindAllOrderByUserID(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, helper.ResponseAPI(false, http.StatusUnauthorized, "unauthorized", gin.H{}))
		return
	}

	//Get All Balance
	balances, err := b.service.FindAllBalanceByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

	c.JSON(http.StatusOK, helper.ResponseAPI(true, http.StatusOK, "success get all balance", balances))
	return
}

func (b *BalanceHandler) FindDetailBalance(c *gin.Context) {
	balanceID := c.Param("id")
	if balanceID == "" {
		c.JSON(http.StatusBadRequest, helper.ResponseAPI(false, http.StatusBadRequest, "bad request", gin.H{}))
		return
	}

	//Get Service
	balance, err := b.service.GetDetailBalanceByID(c, balanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

	c.JSON(http.StatusOK, helper.ResponseAPI(true, http.StatusOK, "success get balance", balance))
	return
}
